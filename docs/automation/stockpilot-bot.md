# StockPilot Bot

StockPilot Bot is a small, repository-native GitHub Actions assistant for issue conversations.

## Commands

Use one of these commands in an issue or pull-request conversation:

- `/stockpilot help` — list supported commands.
- `/stockpilot status` — show the automation status message.
- `/stockpilot triage` — provide a structured bug-triage checklist.

A newly opened issue that contains `/stockpilot` also receives the command response automatically.

## Design and security

The bot runs from `.github/workflows/bot.yml` and uses `actions/github-script@v8`.

The workflow intentionally:

- grants `contents: read` only;
- grants `issues: write` for comments;
- grants `pull-requests: read` for pull-request context;
- does not check out repository code;
- does not execute issue or pull-request text as shell commands;
- does not merge pull requests;
- does not deploy StockPilot;
- does not modify StockPilot application or database data.

## Extending the bot

Future bot commands should remain deterministic and least-privileged. Any command that would write repository files, merge code, publish releases, or access production systems should be implemented as a separate reviewed workflow with narrowly scoped permissions rather than adding broad permissions to this bot.
