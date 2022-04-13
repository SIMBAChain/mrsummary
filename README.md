MRSummary
======================
[![build](https://github.com/SIMBAChain/mrsummary/actions/workflows/build.yaml/badge.svg?branch=main)](https://github.com/SIMBAChain/mrsummary/actions/workflows/build.yaml)  [![release](https://github.com/SIMBAChain/mrsummary/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/SIMBAChain/mrsummary/actions/workflows/release.yaml) 

Simple tool for summarising a GitLab MR in terms of Jira tickets.

[Get the latest release here](https://github.com/SIMBAChain/mrsummary/releases)

### Running

Easiest usage is to pass in the details through environmental variables.

```bash
Usage:
  mrsummary [OPTIONS]

Application Options:
      --jiraurl=     Jira instance URL [$JIRA_URL]
      --jiraemail=   Jira API Token [$JIRA_EMAIL]
      --jiratoken=   Jira API Token [$JIRA_TOKEN]
      --project=     Gitlab Project ID [$CI_MERGE_REQUEST_PROJECT_ID]
      --mergeiid=    Gitlab Merge Internal ID [$CI_MERGE_REQUEST_IID]
      --gitlabtoken=

Help Options:
  -h, --help         Show this help message
```