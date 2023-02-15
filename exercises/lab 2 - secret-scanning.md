## Secret scanning

### Contents

- [Enabling secret scanning](#enabling-secret-scanning)
- [Viewing and managing results](#viewing-and-managing-results)
- [Introducing a test secret](#introducing-a-test-secret)
- [Excluding files from secret scanning](#excluding-files-from-secret-scanning)
- [Managing access to alerts](#managing-access-to-alerts)

### _**Lab 2**_

#### Enabling secret scanning
Secret scanning can be enabled in the settings of an organization or a repository.

1. Go to the repository settings and enable secret scanning in the `Code security and analysis` section.

#### Viewing and managing results
After a few minutes, the `Security` tab in the repository will indicate that there are new security alerts.

- Go to the `Secret scanning` section to view the detected secrets.

For each secret, look at the options to close it and determine which one is most suitable.

#### Introducing a test secret
When developing test cases, you might find that secrets are introduced that cannot be abused when disclosed. Secret scanning will still detect and alert on these secrets.

1. In the GitHub repository file explorer, create a test file that contains a test secret.
    - For example the file `storage-service/src/main/resources/application.dev.properties` with the secrets
        ```
        AWS_ACCESS_KEY_ID="AKIAZBQE345LKPTEAHQD"
        AWS_SECRET_ACCESS_KEY="wt6lVzza0QFx/U33PU8DrkMbnKiu+bv9jheR0h/D"
        ```
2. Determine if the secret is detected when the file is stored.
3. How would you like to manage results from test files?

#### Excluding files from secret scanning
While we can close a detected secret as being used in a test, we can also configure secret scanning to exclude files from being scanned.

1. Create the file `.github/secret_scanning.yml` if it doesn't already exist.
2. Add a list of paths to exclude from secret scanning. You can use [filter patterns](https://docs.github.com/en/free-pro-team@latest/actions/reference/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet) to specify paths.
    ```yaml
    paths-ignore:
        - '**/test/**'
    ```
    **Note**: The characters `*`, `[`, and `!` are special characters in YAML. If you start a pattern with `*`, `[`, or `!`, you must enclose the pattern in quotes.

    Use a pattern to exclude the file `storage-service/src/main/resources/application.dev.properties`

    Merge your changes to `.github/secret_scanning.yml` to your default branch before going to the next step.

    <details>
    <summary>Solution</summary>
    A possible solution is:

    ```yaml
    paths-ignore:
        - '**/test/**'
        - '**/application.dev.properties'
    ```
    </details>

3. Test the pattern by adding another secret or to the file `storage-service/src/main/resources/application.dev.properties`

    For example change the `secretKey` to
    ```
    AWS_SECRET_ACCESS_KEY="6L=yQr6Ivxxj/XG+YdFPdH/xWDcbSV9ch/EjmHCL"
    ```

#### Custom patterns for secret scanning
Secret scanning supports finding other [secret patterns](https://docs.github.com/en/code-security/secret-security/defining-custom-patterns-for-secret-scanning), which are specified by regex patterns and use the Hyperscan library.

1. Add a custom secret pattern by going to the `Code security and analysis` settings and under the header "Custom patterns" click on `New pattern`.
2. Add a custom pattern name, a secret format and test cases.

    For example:
    ```
    Custom pattern name: My secret pattern
    Secret format: my_custom_secret_[a-z0-9]{3}
    Test string: my_custom_secret_123
    ```
 3. Save your pattern and observe the secret scanning alerts page to see if your custom secret pattern has been detected.

#### Managing access to alerts
Due to the nature of secrets, the alerts are only visible to organization and repository administrators and users with `write` or higher access to the repo .
Access to other members and teams can be given in the `Code security and analysis` setting.

**Note:** The member or teams require write privileges before access to alerts can be given.

- In the `Access to alerts` section, add another team member or team to provide access to your repository alerts.


💡**Now that we're familiar with secret scanning, let's head over to the code scanning section, and learn more about it!** 💡
