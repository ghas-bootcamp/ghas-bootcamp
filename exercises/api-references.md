## Additional References

### Code scanning API
The code scanning REST API can be used to retrieve information or modify existing information. Explore the options [here](https://docs.github.com/en/free-pro-team@latest/rest/reference/code-scanning).

### Secret scanning API
The secret scanning API lets you retrieve and update secret scanning alerts from a private repository. Explore the options [here](https://docs.github.com/en/rest/reference/secret-scanning)

### Automating Dependabot

For the Dependabot service, we have both REST and GraphQL endpoints that allow us to configure aspects and retrieve information.
This can both be used for managing individual repositories at scale or export vulnerability information into other systems used to manage this information in your organization. We've included a postman collection in this repository for you to explore.

**Note** The REST API signals a `404` when you client isn't properly authenticated to limit disclosure of private repositories as outlined [here](https://docs.github.com/en/free-pro-team@latest/rest/overview/troubleshooting#why-am-i-getting-a-404-error-on-a-repository-that-exists). This is the same status code that is returned if features are not enabled.

1. Start by creating a *Personal Access Token* with the `repo` scope and store it for use (e.g., in a Postman environment variable).
1. For your repo, determine if vulnerability alerts are enabled.

   **Hints**
   - [Check if vulnerability alerts are enabled for a repository](https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#check-if-vulnerability-alerts-are-enabled-for-a-repository
   )
   - Since this API is currently available for developers preview we need to use the header `Accept: application/vnd.github.dorian-preview+json`
1. Disable and Enable security updates

   **Hints**
   - [Enable automated security fixes](https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#enable-automated-security-fixes)

<details>
<summary>Solutions</summary>

1. Determining if vulnerability alerts are enabled

    ```bash
    curl --location --request GET 'https://api.github.com/repos/<owner>/<repository>/vulnerability-alerts' \
    --header 'Accept: application/vnd.github.dorian-preview+json' \
    --header 'Authorization: Bearer <insert your PAT>'
    ```
1. Disabling and enable security updates

    ```bash
    curl --location --request DELETE 'https://api.github.com/repos/<owner>/<repository>/vulnerability-alerts' \
    --header 'Accept: application/vnd.github.dorian-preview+json' \
    --header 'Authorization: Bearer <insert your PAT>'

    curl --location --request PUT 'https://api.github.com/repos/<owner>/<repository>/vulnerability-alerts' \
    --header 'Accept: application/vnd.github.dorian-preview+json' \
    --header 'Authorization: Bearer <insert your PAT>'
    ```
</details>

Next, we are going to use the GraphQL API to retrieve information on vulnerable dependencies in our repository.

1. Retrieve the `securityVulnerability` objects for your repository.

   If you receive this response, then you need to add the scope `read:org` to your *Personal Access Token*.

   ```json
    {
        "data": {
            "viewer": {
                "organization": null
            }
        }
    }
   ```
   **Hints**
   1. GraphQL is introspective; you can query an object's schema with

        ```graphql
        query {
            __type(name: "SecurityVulnerability") {
                name
                kind
                description
                fields {
                    name
                }
            }
        }
        ```
   1. A [SecurityVulnerability](https://docs.github.com/en/free-pro-team@latest/graphql/reference/objects#securityvulnerability) object can be accessed via the [RepositoryVulnerabilityAlert](https://docs.github.com/en/free-pro-team@latest/graphql/reference/objects#repositoryvulnerabilityalert) object in a [Repository](https://docs.github.com/en/free-pro-team@latest/graphql/reference/objects#repository) object, which itself resides in an [Organization](https://docs.github.com/en/free-pro-team@latest/graphql/reference/objects#organization) object.


<details>
<summary>Solution</summary>

```graphql
query VulnerabilityAlerts($org: String!, $repo: String!){
  viewer {
    organization(login: $org) {
      repository(name: $repo) {
        name
        vulnerabilityAlerts(first: 10) {
          nodes {
            securityVulnerability {
              advisory {
                ghsaId
                description
              }
              package {
                name
                ecosystem
              }
              severity
              firstPatchedVersion {
                identifier
              }
              vulnerableVersionRange
            }
          }
        }
      }
    }
  }
}
```

```bash
curl --location --request POST 'https://api.github.com/graphql' \
--header 'Authorization: Bearer <insert your PAT>' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"query VulnerabilityAlerts($org: String!, $repo: String!){\n  viewer {\n    organization(login: $org) {\n      repository(name: $repo) {\n        name\n        vulnerabilityAlerts(first: 10) {\n          nodes {\n            securityVulnerability {\n              advisory {\n                ghsaId\n                description\n              }\n              package {\n                name\n                ecosystem\n              }\n              severity\n              firstPatchedVersion {\n                identifier\n              }\n              vulnerableVersionRange\n            }\n          }\n        }\n      }\n    }\n  }\n}","variables":{"org":"<org-name>","repo":"<repository-name>"}}'
```

</details>
