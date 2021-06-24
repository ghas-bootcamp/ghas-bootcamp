## Dependabot

### Contents

- [Enabling Dependabot alerts](#enabling-dependabot-alerts)
- [Reviewing the dependency graph](#reviewing-the-dependency-graph)
- [Viewing and managing results](#viewing-and-managing-results)
- [Enabling Dependabot security updates](#enabling-dependabot-security-updates)
- [Configuring Dependabot security updates](#configuring-dependabot-security-updates)

### _**Practical Exercise 1b**_

#### Enabling Dependabot alerts
Dependabot can be enabled in the settings of an organization or a repository.

1. Go to the repository settings and enable Dependabot alerts in the *Security & analysis* section.

#### Reviewing the dependency graph
Dependabot uses the dependency graph to determine which dependencies are used by your project.

1. Verify in the dependency graph that it found dependencies for:
    - The frontend service.
    - The authentication service.
    - The gallery service.
    - The storage service.

The dependency graph can be access in the `Insights` tab in your repository.

#### Viewing and managing results

After a few minutes, the security tab in the repository will indicate that there are new security alerts.

**Note**: If this not the case, we can trigger an analysis by updating `authn-service/requirements.txt`

1. Go to the Dependabot alert section to view the detected dependency issues.

For each dependency alert, we have the option to create a security update or to dismiss the alert with a reason.

2. For one of the alerts create a dependency security update. If Dependabot can update the dependency automatically, it will create a PR.

3. For one of the alerts dimiss the alert.

#### Enabling Dependabot security updates

Dependabot can automatically create PRs to upgrade vulnerable dependencies to non-vulnerable versions.

1. Go to the repository settings and enable Dependabot security updates in the *Security & analysis* section.

After a few minutes multiple PRs will be created that will upgrade vulnerable dependencies.

#### Configuring Dependabot security updates

You can enable Dependabot [*version updates*](https://docs.github.com/en/code-security/supply-chain-security/keeping-your-dependencies-updated-automatically/enabling-and-disabling-version-updates) by checking in a dependabot.yml into your repository's `.github` directory. Dependabot security updates uses this configuration as well. To successfully integrate the security updates into the SDLC, it is possible to configure various aspects such as:

- When security PRs are created.
- What labels are assigned to enable filtering options.
- Who is assigned to the PR and who should review it.
- Specify which dependency are updated and how they are updated.

Create the file `.github/dependabot.yml` in your repository and configure the `pip` dependency manager to:
  1. Look for dependency information in the directory `authn-service`.

  2. Schedule daily security updates.

  3. Prefix the commit message with the package manager `pip`.

  4. Assign the PR to yourself and a person from your workshop team as a reviewer.

  5. Add the custom label `triage-required` to enable filtering of the PRs (Make sure the label exists by adding it to `https://github.com/<owner>/<your repo>/labels`).

  6. Verify your changes by adding a [vulnerable dependency](https://github.com/advisories?query=severity%3Ahigh+ecosystem%3Apip) to `auth-service/requirements.txt`. For example:

    ```requirements.txt
    ...
    django==2.1.0
    ```

How would you know if the configuration cannot be satisfied?

1. Add a non-existing label to the configuration.

2. Trigger a new dependabot security update by adding a vulnerable dependency to one of the projects
   For example, we can add the dependency `django-two-factor-auth==1.11` to `auth-service/requirements.txt`

3. Look at the created PR to determine if the configuration has been satisfied.

<details>
<summary>Solution</summary>

```yaml
version: 2
updates:
  - package-ecosystem: "pip"
    directory: "/authn-service"
    schedule:
      interval: "daily"
    labels:
      - "triage-required"
    assignees:
      - "<github handle>"
    reviewers:
      - "<github handle>"
    commit-message:
      prefix: "pip"
```
</details>
    
#### Working with Dependency Review
    
If a Pull Request has dependency changes, you can [review](https://docs.github.com/en/github/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/reviewing-dependency-changes-in-a-pull-request) them and see if there are known vulnerabilities with the dependency changes.
    
   1. Add a vulnerable dependency to `auth-service/requirements.txt` and commit to a new branch. For example, here's a vulnerable dependency:

    ```requirements.txt
    ...
    django-piston==0.2.0
    ```
   2. Create a Pull Request, and click on `Files changed`. 
   3. Click on the `Display the rich diff` button on the `requirements.txt` file to review dependency changes.

💡**Now that we're familiar with Dependabot, let's head over to the code scanning section, and learn more about it! [Click here](code-scanning.md).** 💡
