import contextlib
import jinja2
import networkx as nx
import os
import yaml

from collections import defaultdict
from faker import Faker
from typing import *

try:
    from django.contrib.auth.hashers import make_password as hash_password
except ImportError:
    hash_password = lambda x: x

from .factories import FactoryInterface
from .types import Identifier
from .utils import normalize_identifiers, random_model, random_models


MULTI_CLASS_FILENAMES = {'fixtures.yml', 'fixtures.yaml'}


class FixturesLoader:
    """
    The factory "driver" class. Does most of the hard work of loading fixtures,
    leaving the responsibility of model instantiation up to the factory class
    passed in.

    :param factory: An instance of the concrete factory to use for creating models
    :param fixture_dirs: A list of directory paths to load fixtures templates from
    :param env: An optional jinja environment (the default one will include
                faker as a template global, but if you want to customize its
                tags/filters/etc, then you need to create an env yourself - the
                correct loader will be set automatically for you)
    """

    def __init__(self,
                 factory: FactoryInterface,
                 fixture_dirs: List[str],
                 env: Optional[jinja2.Environment] = None):
        self.env = self._ensure_env(env)
        """The Jinja Environment used for rendering the yaml template files."""

        factory.loader = self
        self.factory = factory
        """The factory instance."""

        self.fixture_dirs = fixture_dirs
        """A list of directories where fixture files should be loaded from."""

        self.relationships = {}
        """A dict keyed by model name where values are a list of related model names."""

        self.model_fixtures = defaultdict(dict)
        """A dict of models names to their semi-processed data from the yaml files."""

        self._file_cache = {}
        self._data_cache = defaultdict(dict)
        self._loaded = False

    def create_all(self, progress_callback: Optional[callable] = None) -> Dict[str, object]:
        """
        Creates all the models discovered from fixture files in :attr:`fixtures_dir`.

        :param progress_callback: An optional function to track progress. It must take three
                               parameters:
                                - an :class:`Identifier`
                                - the model instance
                                - and a boolean specifying whether the model was created
        :return: A dictionary keyed by identifier where the values are model instances.
        """
        if not self._loaded:
            self._load_data()

        # build up a directed acyclic graph to determine the model instantiation order
        dag = nx.DiGraph()
        for model_class_name, dependencies in self.relationships.items():
            dag.add_node(model_class_name)
            for dep in dependencies:
                dag.add_edge(model_class_name, dep)

        try:
            creation_order = reversed(list(nx.topological_sort(dag)))
        except nx.NetworkXUnfeasible:
            raise Exception('Circular dependency detected between models: ' +
                            ', '.join('{a} -> {b}'.format(a=a, b=b)
                                      for a, b in nx.find_cycle(dag)))

        # create or update the models in the determined order
        rv = {}
        for model_class_name in creation_order:
            for identifier_key, data in self.model_fixtures[model_class_name].items():
                identifier = Identifier(model_class_name, identifier_key)
                data = self.factory.maybe_convert_values(identifier, data)
                self._data_cache[model_class_name][identifier_key] = data

                model_instance, created = self.factory.create_or_update(identifier, data)
                if progress_callback:
                    progress_callback(identifier, model_instance, created)
                rv[identifier_key] = model_instance

        self.factory.commit()
        return rv

    def convert_identifiers(self, identifiers: Union[Identifier, List[Identifier]]):
        """
        Convert an individual :class:`Identifier` to a model instance,
        or a list of Identifiers to a list of model instances.
        """
        if not identifiers:
            return identifiers

        def _create_or_update(identifier):
            data = self._data_cache[identifier.class_name][identifier.key]
            return self.factory.create_or_update(identifier, data)[0]

        if isinstance(identifiers, Identifier):
            return _create_or_update(identifiers)
        elif isinstance(identifiers, list) and isinstance(identifiers[0], Identifier):
            return [_create_or_update(identifier) for identifier in identifiers]
        else:
            raise TypeError('`identifiers` must be an Identifier or list of Identifiers.')

    def _load_data(self):
        """
        Load all fixtures from :attr:`fixtures_dir`
        """
        filepaths = []
        model_identifiers = defaultdict(list)

        # attempt to load fixture files from given directories (first pass)
        # for each valid model fixture file, read it into the cache and get the
        # list of identifier keys from it
        for fixtures_dir in self.fixture_dirs:
            for filename in os.listdir(fixtures_dir):
                filepath = os.path.join(fixtures_dir, filename)
                file_ext = filename[filename.find('.')+1:]

                # make sure it's a valid fixture file
                if os.path.isfile(filepath) and file_ext in {'yml', 'yaml'}:
                    filepaths.append(filepath)
                    with open(filepath) as f:
                        self._file_cache[filepath] = f.read()

                    # preload to determine identifier keys
                    with self._preloading_env() as env:
                        rendered_yaml = env.get_template(filepath).render()
                        data = yaml.load(rendered_yaml, Loader=yaml.FullLoader)
                        if data:
                            if filename in MULTI_CLASS_FILENAMES:
                                for class_name in data:
                                    model_identifiers[class_name] = list(
                                        data[class_name].keys())
                            else:
                                class_name = filename[:filename.rfind('.')]
                                model_identifiers[class_name] = list(data.keys())

        # second pass where we can render the jinja templates with knowledge of all
        # the model identifier keys (allows random_model and random_models to work)
        for filepath in filepaths:
            self._load_from_yaml(filepath, model_identifiers)

        self._loaded = True

    def _load_from_yaml(self, filepath: str, model_identifiers: Dict[str, List[str]]):
        """
        Load fixtures from the given filename
        """
        rendered_yaml = self.env.get_template(filepath).render(
            model_identifiers=model_identifiers)
        data = yaml.load(rendered_yaml, Loader=yaml.FullLoader)

        identifier_data = {}
        filename = os.path.basename(filepath)
        if filename in MULTI_CLASS_FILENAMES:
            for class_name in data:
                d, self.relationships[class_name] = self._post_process_yaml_data(
                    data[class_name], self.factory.get_relationships(class_name))
                identifier_data[class_name] = d
        else:
            class_name = filename[:filename.rfind('.')]
            d, self.relationships[class_name] = self._post_process_yaml_data(
                data, self.factory.get_relationships(class_name))
            identifier_data[class_name] = d

        for class_name, d in identifier_data.items():
            for identifier_key, instance_data in d.items():
                self.model_fixtures[class_name][identifier_key] = instance_data

    def _post_process_yaml_data(self,
                                fixture_data: Dict[str, Dict[str, Any]],
                                relationship_columns: Set[str],
                                ) -> Tuple[Dict[str, Dict[str, Any]], List[str]]:
        """
        Convert and normalize identifier strings to Identifiers, as well as determine
        class relationships.
        """
        rv = {}
        relationships = set()
        if not fixture_data:
            return rv, []

        for identifier_id, data in fixture_data.items():
            new_data = {}
            for col_name, value in data.items():
                if col_name not in relationship_columns:
                    new_data[col_name] = value
                    continue

                identifiers = normalize_identifiers(value)
                if identifiers:
                    relationships.add(identifiers[0].class_name)

                if isinstance(value, str) and len(identifiers) <= 1:
                    new_data[col_name] = identifiers[0] if identifiers else None
                else:
                    new_data[col_name] = identifiers

            rv[identifier_id] = new_data
        return rv, list(relationships)

    def _ensure_env(self, env: Union[jinja2.Environment, None]):
        """
        Make sure the jinja environment is minimally configured.
        """
        if not env:
            env = jinja2.Environment()
        if not env.loader:
            env.loader = jinja2.FunctionLoader(lambda path: self._file_cache[path])

        if 'faker' not in env.globals:
            faker = Faker()
            faker.seed_instance(1234)
            env.globals['faker'] = faker

        env.globals.setdefault('hash_password', hash_password)
        env.globals.setdefault('random_model', jinja2.contextfunction(random_model))
        env.globals.setdefault('random_models', jinja2.contextfunction(random_models))

        return env

    @contextlib.contextmanager
    def _preloading_env(self):
        """
        A "stripped" jinja environment.
        """
        ctx = self.env.globals
        try:
            ctx['random_model'] = lambda *a, **kw: None
            ctx['random_models'] = lambda *a, **kw: None
            yield self.env
        finally:
            ctx['random_model'] = jinja2.contextfunction(random_model)
            ctx['random_models'] = jinja2.contextfunction(random_models)
