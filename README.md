# gosmarty
A Smarty template engine interpreter written in Go.

This project aims to provide a Go-native implementation of the popular PHP template engine, Smarty.

## Features
**In development**

### Implemented Features

| Feature                | Syntax Example                                       | Status      |
| ---------------------- | ---------------------------------------------------- | ----------- |
| Variable Rendering     | `{$name}`                                            | ✅ |
| Field Access           | `{$user.name}`                                       | ❌ |
| Array Access           |  `{$users[0].name}`                                   | ❌ |
| Variable Modifiers     | `{$title\|upper\|escape}`                            | ✅ |
| If/Else Statements     | `{if $isLoggedIn}Welcome!{else}Please log in.{/if}`  | ✅ |
|                        | `{if $num > 5}true{else}false{/if}`                  | ❌ |
| Comments               | `{* This is a comment *}`                            | ✅ |

### Roadmap

Our goal is to achieve full compatibility with the PHP Smarty engine. The following features are planned for future releases.

| Phase | Feature                               | Description                                                 |
| ----- | ------------------------------------- | ----------------------------------------------------------- |
| 1     | **Core Functionality**                | `{foreach}`, `{assign}`, `{elseif}`                         |
| 2     | **Template Structuring**              | `{include}`, `{extends}`, `{block}`                         |
| 3     | **Compatibility**                     | Essential built-in functions (`{strip}`, `{literal}`, etc.) |
| 4     | **Advanced Features**                 | Caching, Plugin System                                      |
