# Alfred Workflow Object Schema (Reverse-Engineered)

Alfred does not publish an official schema for the objects inside a
workflow's `info.plist` — not for Script Filter, not for Keyword, not for
Universal Action Trigger, not for anything. The [Alfred help
pages](https://www.alfredapp.com/help/workflows/) only document the GUI;
the raw plist format is undocumented and has to be learned by building an
object in Alfred and inspecting what it writes. This file records that
knowledge so it doesn't have to be rediscovered.

This is why `workflow/info.plist` says "edit it directly (there's no
builder)" in [DEVELOPING.md](../DEVELOPING.md) — there is no official
builder to point at.

## How this reference was generated

1. In Alfred Preferences → Workflows, create a blank workflow and add one
   of every object from each category in the `+` menu (Triggers, Inputs,
   Outputs, Actions, Utilities, User Interface, Automation). Leave every
   field at its default — the goal is the object's config *keys*, not
   meaningful values.
2. Right-click the workflow → **Export Workflow…**.
3. Unzip the `.alfredworkflow` (it's a zip archive) and pretty-print its
   `info.plist`: `plutil -convert xml1 -o pretty.plist info.plist`.

The lookup table below (config key names only, not values) was produced by
exactly that process against a scratch workflow with one of every object
type. If a config key's meaning here turns out to be wrong, redo the
export and inspect the real plist directly — that always out-ranks any
prose description below, including this one.

## General structure

Every workflow's `info.plist` is a dict with (at minimum):

- `objects` — an array of objects, each `{config: {...}, type: "alfred.workflow.…", uid: "<UUID>", version: <int>}`. Order in the array is cosmetic only; Alfred re-sorts it on every save. Never rely on array position — always match by `uid`.
- `connections` — a dict keyed by a source object's `uid`, each value an array of `{destinationuid, modifiers, modifiersubtext, vitoclose}` (one entry per outgoing wire, `modifiers` distinguishing plain vs. modifier-key connections).
- `uidata` — a dict keyed by `uid`, each value `{xpos, ypos}`: purely the canvas layout, never functional.
- `bundleid`, `name`, `description`, `createdby`, `webaddress`, `readme`, `version`, `disabled`, `userconfigurationconfig`, `variablesdontexport` — workflow-level metadata, not object-level.

A `uid` is a stable identifier for that object for as long as the workflow exists — reconnecting objects doesn't change it, but deleting and re-adding one does. When hand-authoring a plist (as this project does), invent your own UUIDs; Alfred accepts any well-formed UUID string.

## Object types this project uses

### `alfred.workflow.input.scriptfilter`

The object behind both `cit` (keyword) and the selection chooser. Key config fields:

| Key | Meaning |
|---|---|
| `keyword` | The typed keyword, or `""` for a keyword-less node driven only by an incoming connection (e.g. from a Universal Action Trigger). |
| `alfredfiltersresults` | If `true`, Alfred fuzzy-filters the returned items' `title`/`match` against the incoming query string before displaying them. **This project sets it `false`** on every keyword-less node, because the incoming query is arbitrary selected text, not something meant to fuzzy-match static item titles like "Check"/"Reveal"/"Clean" — leaving it `true` here silently filters every item away (see the gotcha below). |
| `argumenttreatemptyqueryasnil` | Whether an empty query is passed as `nil` vs `""` to the script. |
| `title` / `subtext` | The static label/placeholder shown before the script runs. Distinct from any individual result item's own `title`/`subtitle`. |
| `runningsubtext` | Shown while the script is executing (e.g. "Reading the clipboard…"). |
| `script` / `scriptargtype` | The command line run, and whether `$1`/`argv` or `{query}` substitution is used. |

**Gotcha this project actually hit**: a chooser like this one (`list selection "$1"`, keyword-less, `subtext` "Check, reveal, or clean the selected text") and the node one step downstream that executes the chosen action (`run "$1"`, also keyword-less, but with blank `title`/`subtext`, `runningsubtext` "Running…") are easy to visually confuse in the canvas — both display as unlabelled-looking boxes. A Universal Action Trigger wired directly to the *executor* node instead of the *chooser* node skips the Check/Reveal/Clean selection screen entirely: the raw selected text arrives as the executor's argument, matches none of its recognized actions, and it falls back to "Nothing to do" for every single selection. Identify the correct target by `subtext`, not by title (both nodes show the same title, "Clean Invisible Text").

### `alfred.workflow.action.script`

A plain "Run Script" step with no result list of its own (used here for `copy-report`). Config: `script`, `scriptargtype`, `scriptfile`, `escaping`, `concurrently`.

### `alfred.workflow.trigger.universalaction`

```xml
<dict>
    <key>config</key>
    <dict>
        <key>acceptsfiles</key><false/>
        <key>acceptsmulti</key><integer>0</integer>
        <key>acceptstext</key><true/>
        <key>acceptsurls</key><false/>
        <key>name</key><string>Clean Invisible Text</string>
    </dict>
    <key>type</key><string>alfred.workflow.trigger.universalaction</string>
    <key>uid</key><string>…</string>
    <key>version</key><integer>1</integer>
</dict>
```

Entirely static — no machine-specific state, no registration step baked into the plist. `name` is the label shown in Alfred's Universal Actions list. `acceptsmulti` (an integer, not a bool) is `0` for "single selection only" — this project uses `0` because [`internal/action.Request`](../internal/action/action.go) takes one string, not an array; a workflow meant to batch-process multiple selections would need `1` here and array-handling downstream. `acceptstext`/`acceptsfiles`/`acceptsurls` gate which selection types Alfred offers this action for.

Contrary to what this project assumed for a long time (see `git log` on `README.md`/`DEVELOPING.md` before this was corrected), this object **can** be committed to `workflow/info.plist` and works immediately on import — no one-time manual setup is required. The earlier "isn't something this project can generate reproducibly" claim was never actually tested; it turned out to be as reproducible as every other object type once someone exported a workflow and looked.

## Other object types (lookup table)

The "UI name" column is the *official* name, confirmed against Alfred's own
help pages (linked per category below) — not a guess. Config keys are
still from inspecting the exported plist per the process above, since
Alfred's help pages document usage, not the raw plist field names.

| Type | Config keys | UI name |
|---|---|---|
| `alfred.workflow.action.actioninalfred` | `path`, `type` | [Action in Alfred](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.applescript` | `applescript`, `cachescript` | [Run NSAppleScript](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.browseinalfred` | `path`, `sortBy`, `sortDirection`, `sortFoldersAtTop`, `sortOverride`, `stackBrowserView` | [Browse in Alfred](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.browseinterminal` | `path` | [Browse in Terminal](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.buffer` | `addfilestobuffer`, `clearbuffer`, `outputtype` | [File Buffer](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.itunescommand` | `command` | [Music Command](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.launchfiles` | `paths`, `toggle` | [Launch Apps / Files](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.openfile` | `openwith`, `sourcefile` | [Open File](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.openurl` | `browser`, `skipqueryencode`, `skipvarencode`, `spaces`, `url` | [Open URL](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.revealfile` | `path` | [Reveal File in Finder](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.systemcommand` | `command`, `confirm` | [System Command](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.systemwebsearch` | `browser`, `searcher` | [Default Web Search](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.action.terminalcommand` | `escaping`, `script` | [Terminal Command](https://www.alfredapp.com/help/workflows/actions/) |
| `alfred.workflow.automation.runshortcut` | `inputmode`, `outputmode`, `shortcut` | [Run Shortcut](https://www.alfredapp.com/help/workflows/automations/) (macOS 12+ Shortcuts app) |
| `alfred.workflow.automation.task` | *(none)* | [Automation Task](https://www.alfredapp.com/help/workflows/automations/) |
| `alfred.workflow.input.dictionaryfilter` | `keyword`, `language`, `showallwords`, `subtext`, `title` | [Dictionary Lookup](https://www.alfredapp.com/help/workflows/inputs/) |
| `alfred.workflow.input.filefilter` | `anchorfields`, `argumenttrimmode`, `argumenttype`, `daterange`, `fields`, `includesystem`, `limit`, `runningsubtext`, `scopes`, `sortmode`, `subtext`, `title`, `types`, `withspace` | [File Filter](https://www.alfredapp.com/help/workflows/inputs/) |
| `alfred.workflow.input.keyword` | `argumenttype`, `subtext`, `text`, `withspace` | [Keyword](https://www.alfredapp.com/help/workflows/inputs/) |
| `alfred.workflow.input.listfilter` | `argumenttrimmode`, `argumenttype`, `fixedorder`, `items`, `matchmode`, `runningsubtext`, `subtext`, `title`, `withspace` | [List Filter](https://www.alfredapp.com/help/workflows/inputs/) |
| `alfred.workflow.input.runningappsfilter` | `argumenttype`, `keyword`, `outputprefix`, `outputtype`, `subtext`, `title`, `withspace` | [Running Apps Filter](https://www.alfredapp.com/help/workflows/inputs/) |
| `alfred.workflow.output.callexternaltrigger` | `externaltriggerid`, `passinputasargument`, `passvariables`, `workflowbundleid` | [Call External Trigger](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.clipboard` | `autopaste`, `clipboardtext`, `ignoredynamicplaceholders`, `transient` | [Copy to Clipboard](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.dispatchkeycombo` | `count`, `keychar`, `keycode`, `keymod`, `overridewithargument` | [Dispatch Key Combo](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.largetype` | `alignment`, `backgroundcolor`, `fadespeed`, `fillmode`, `font`, `ignoredynamicplaceholders`, `largetypetext`, `textcolor`, `wrapat` | [Large Type](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.notification` | `lastpathcomponent`, `onlyshowifquerypopulated`, `removeextension`, `text`, `title` | [Post Notification](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.playsound` | `soundname`, `systemsound` | [Play Sound](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.speak` | `text`, `usevoiceover` | [Speak](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.output.writefile` | `adduuid`, `allowemptyfiles`, `createintermediatefolders`, `filename`, `filetext`, `ignoredynamicplaceholders`, `relativepathmode`, `type` | [Write File](https://www.alfredapp.com/help/workflows/outputs/) |
| `alfred.workflow.trigger.action` | `acceptsmulti` | [File Action](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.trigger.contact` | *(none)* | [Contact Action](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.trigger.external` | `availableviaurlhandler` | [External](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.trigger.fallback` | *(none)* | [Fallback Search](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.trigger.hotkey` | `action`, `argument`, `focusedappvariable`, `focusedappvariablename`, `hotkey`, `hotmod`, `leftcursor`, `modsmode`, `relatedAppsMode` | [Hotkey](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.trigger.remote` | `argumenttype`, `workflowonly` | [Remote](https://www.alfredapp.com/help/workflows/triggers/) (Alfred Remote app) |
| `alfred.workflow.trigger.snippet` | `focusedappvariable`, `focusedappvariablename` | [Snippet](https://www.alfredapp.com/help/workflows/triggers/) |
| `alfred.workflow.userinterface.grid` | `columncount`, `filterable`, `fixedorder`, `imageaspect`, `inputfile`, `inputtype`, `loadingtext`, `showsubtitles`, `showtitles`, `subtitlesinfooter`, `titlesinfooter` | [Grid View](https://www.alfredapp.com/help/workflows/user-interface/) |
| `alfred.workflow.userinterface.image` | `imageresizemode`, `stackview` | [Image View](https://www.alfredapp.com/help/workflows/user-interface/) |
| `alfred.workflow.userinterface.pdf` | `displaymode`, `stackview` | [PDF View](https://www.alfredapp.com/help/workflows/user-interface/) |
| `alfred.workflow.userinterface.text` | `behaviour`, `fontmode`, `fontsizing`, `footertext`, `inputfile`, `inputtype`, `loadingtext`, `outputmode`, `scriptinput`, `spellchecking`, `stackview` | [Text View](https://www.alfredapp.com/help/workflows/user-interface/) |
| `alfred.workflow.utility.argument` | `argument`, `passthroughargument` | [Arguments and Variables](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.conditional` | `conditions`, `elselabel`, `hideelse` | [Conditional](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.debug` | `argument`, `cleardebuggertext`, `processoutputs` | [Debug](https://www.alfredapp.com/help/workflows/utilities/debug/) |
| `alfred.workflow.utility.delay` | `seconds` | [Delay](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.dialog` | `button1`, `button2`, `button3`, `description`, `title` | [Dialog Conditional](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.expression` | `expression` | [Expression](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.file` | `fileutivariablename`, `outputfileuti` | [File Conditional](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.filter` | `inputstring`, `matchcasesensitive`, `matchmode`, `matchstring` | [Filter](https://www.alfredapp.com/help/workflows/utilities/) (legacy) |
| `alfred.workflow.utility.hidealfred` | `unstackview` | [Hide Alfred](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.joinargs` | `delimiter` | [Join Args](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.json` | `json` | [JSON Configuration](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.junction` | *(none)* | [Junction](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.random` | `type` | [Random](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.replace` | `matchmode`, `matchstring`, `replacestring` | [Replace](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.showalfred` | `argument`, `leftcursor` | [Show Alfred](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.split` | `delimiter`, `discardemptyarguments`, `outputas`, `trimarguments`, `variableprefix` | [Split Arg](https://www.alfredapp.com/help/workflows/utilities/) |
| `alfred.workflow.utility.transform` | `type` | [Transform](https://www.alfredapp.com/help/workflows/utilities/) |

Not listed: `alfred.workflow.input.scriptfilter`, `alfred.workflow.action.script`, and `alfred.workflow.trigger.universalaction` — documented above instead, since this project actually depends on their exact behavior.
