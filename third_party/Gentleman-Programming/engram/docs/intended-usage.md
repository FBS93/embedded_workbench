# Intended Usage

<- [Back to README](../README.md)

---

Engram is persistent memory for AI coding agents. It saves decisions, discoveries, bug fixes, and context across sessions -- automatically. This page explains how it's meant to be used. Not the API, not the architecture -- just the mental model.

---

## If You Installed via gentle-ai -- You're Done

Engram is already configured. Your AI agent saves and retrieves memories on its own.

You never need to configure it, inspect it, or interact with it directly. If engram is working correctly, you won't even notice it's there. That's the point.

---

## If You're Integrating Engram Into Your Agent

Install the plugin (Claude Code, OpenCode) or use the preset. The memory protocol is already included -- it tells your agent what to save, when to search, and how to manage sessions.

You don't need to teach the agent anything. The protocol handles proactive saves, session summaries, and memory search automatically.

That's it. There is no step two.

---

## The Golden Rule

Engram is infrastructure. Like a good database, you set it up once and forget about it.

If you're thinking about engram while working, something went wrong. It should be invisible.

---

## Quick Reference

| Do | Don't |
|----|-------|
| Install via gentle-ai, plugin, or preset | Manually inspect or edit engram's storage |
| Trust that your agent is saving context for you | Try to manage what the agent saves |
| Just start coding -- memory happens in the background | Build custom memory protocols -- use the one that ships with the plugin |
| Re-run the installer to update | Worry about what's being saved or when |
