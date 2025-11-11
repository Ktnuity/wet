# wet [<img src="https://img.shields.io/github/license/ktnuity/wet">](https://github.com/Ktnuity/wet/blob/master/LICENSE)
**wet** is a stack-based scripting language used to set-up projects or codebases immediately after a vcs clone.

It's:
- At some points a pre-build.
- At other parts a post-build.
- You may even call it a build-system.

# What is `wet` really?
Wet is a collection of tasks you'd usually see run in a project to initialize tools and set it up for active development. It could be downloading a file, unzipping that file, placing some of its content in different locations, making everything simple for a new user by doing all the small work for you.

It's sort of there in the name. `wet` is the liquid as it's poured from the kettle into every cup of the dozens of people visiting you. You shouldn't need to do it yourself, `wet` is the servant you've hired to do it for you, and it does exactly as you've told it to.

# Features?
This project is WIP so what follows below is mostly a plan for what will become. And more will of course follow once it's completed.

<details>
    <summary><code>Language</code>: Information about the language that powers <ins>wet</ins></summary>

## Arithmetics
`wet` would support arithmetics similar to that of `porth`. The way this would work is by utilizing [Reverse Polish Notation](https://en.wikipedia.org/wiki/Reverse_Polish_notation). So instead of `x + y` as we humans normally do it, we'd use `x y +`. `wet` would support all of these:
- Arithmetics: `+`, `-`, `/`, `*`, `%`, `++`, `--`
- Bitwise: `&`, `|`, `^`, `~`
- Logical: `=`, `!=`, `<`, `<=`, `>`, `>=`, `!`
- Boolean: `&&`, `||`

## Output
- Strings: `puts`
- Numbers: `.`

## Conditionals & Branching
The following branches and conditionals would exist:
- Loops: `while`, `until`
- Branching: `if`, `unless`

## Types
The following types would be present:
- String: `"hello, world!"`
- Numbers: `123` (int32), `13.37` (float32)
- Bool: `true`, `false`, `1`, `0`
- Resource Location: `/<filepath>` (file relative to `git` root), `./<filepath>` (file relative to current work dir of script), `:<token>` (file stored as a virtual token)

Resource Locations would use `\ ` to excape a space, and `\\` to escape a backslash.

## Operands
- `dup`: pops top value, pushes it back twice.
  - ```
       | A
    -> | A A
    ```
- `over`: pushes 2nd top-most value.
  - ```
       | A B
    -> | A B A
    ```
- `swap`: swap 2 top-most values.
  - ```
       | A B
    -> | B A
    ```
- `2dup`: duplicates the 2 top-most values.
  - ```
       | A B
    -> | A B A B
    ```
- `2swap`: swaps top values 1 & 2 with 3 & 4.
  - ```
       | A B C D
    -> | C D A B
    ```
- `drop`: drops the top most value.
  - ```
       | A B
    -> | A
    ```
- `nop`: does absolutely nothing (sort of like a hardware wait cycle).
  - ```
       | A
    -> | A
    ```

## Memory
- `<value> "<name>" store`: store `<value>` as `@<name>`.
  - consumes `<value>` and `"<name>"`.
- `"<name>" load`: load from `@<name>`.
  - consumes `"<name>"`
  - pushes `<value> true` if successful and type matches, `false` otherwise.
- `"<name>"` is local-scope by default.
  - prefix `<name>` with `.` (`".<name>"`) for global scope.

## Macros
Macros are snippets of code which may be replicated multiple times across the program.<br>
You define a macro by writing `macro <name>`, and you end it as usual by using `end`.

Example:
```
macro my_print
    dup tostring "\n" + puts
end

34 my_print
35 my_print
+ my_print
drop
```

`my_print` here is a macro which prints any stack value onto the screen without consuming it.<br>
When this code is executed, the pre-processor removes all macros and turn it into:
```
34 dup tostring "\n" + puts
35 dup tostring "\n" + puts
+ dup tostring "\n" + puts
drop
```

As you can see, macros don't exist at runtime. Because of this, you cannot use macros for recursion.

## Procedures
Procedures are just like macros, but they persist at runtime, and are called whenever a name is encountered.
You define a procedure by writing `proc <name>`, and you end it with `end`.

Example:
```
proc my_print
    dup tostring "\n" + puts
end

34 my_print
35 my_print
+ my_print
drop
```

Unlike with macros, this code looks the same at runtime.

Let's step through parts of the code together.
1. Instruction is `proc`. `proc` marks a procedure, so we skip to just after `end`.
2. Instruction is `34`. 34 is pushed to stack.
3. Instruction is `my_print`. IP after `my_print` (IP of `35`) is pushed to call stack. IP is changed to inside `my_print` body.
4. Instruction is `dup`. Pops `34` from stack and push it back twice.
5. Instruction is `tostring`. Pops `34` from stack. Turns it from `int` to `string`. Push `"34"` to stack.
6. Instruction is `"\n"`. Push `"\n"` to stack.
7. Instruction is `+`. Pop `"34"` and `"\n"` from stack, concatenate them, and push back `"34\n"`.
8. Instruction is `puts`. Pop `"34\n"` from stack, consume it, and print it.
9. Instruction is `end`. `end` has an `EndMode` which now is set to `EndModeProc`. Pops top value from call stack and change IP to it.
10. Instruction is `35`. 35 is pushed to stack.

You get the idea at this point.

Procedures have the benefit that because they exist in the code at runtime, they support recursion. Check [`demo/proc/init.wet`](./demo/proc/init.wet) for an implementation of Fibonacci procedure in wet.

**Note:** Prefer using procedures when the body is fairly large, and macros when the body is fairly small.

### Returning

Returning uses 3 keywords.
- `ret` immeditately returns from the procedure.
- `<drop> dret` drops the top `<drop>` (count) values from stack before returning.
- `<drop> <keep> iret` stores top `<keep>` (count) values from the stack, drops the next `<drop>` (count) values, then restore the previous `<keep>` (count) values, then returns.

</details>

<details>
    <summary><code>Tools</code>: Tools and commants that's built into <ins>wet</ins>'s build system</summary>

## Tokens vs Files
When a task downloads to a file, it would need an absolute file in relation to your `.git` directory location. However, in order to allow for temporary paths, you will be able to use tokens. To create a token file, you write `:<token>` instead of `<file>`. The prefix of `:` indicates a token being used. On the back-end, this would use a `~/.wet/` for storage, or `./.wet` if manually created next to `./.git`.

## Commands
- `<url> <dst> download`
  - `<url>` expects a `string` containing a remote URL.
  - `<dst>` expects any resource location.
  - downloads the file/resource at `<url>` into `<dst>`
    - both `<url>` and `<dst>` is consumed.
    - pushes `true` if successful, `false` if unsuccessful.
- `<src> readfile`
  - `<src>` expects any resource location.
  - loads the file into memory.
    - `<src>` is consumed.
    - pushes `<data> true` if successful, `false` if unsuccessful.
      - `<data>` is a string containing file content.
      - fails if file is binary.
- `<src> <dst> move`
  - `<src>` and `<dst>` both expects any resource location.
  - moves `<src>` to `<dst>`
    - both `<src>` and `<dst>` is consumed.
    - pushes `true` if successful, `false` if unsuccessful.
- `<src> <dst> copy`
  - `<src>` and `<dst>` both expects any resource location.
  - copies `<src>` to `<dst>`
    - both `<src>` and `<dst>` is consumed.
    - pushes `true` if successful, `false` if `<dst>` already exist or if copy was unsuccessful.
- `<res> exist`
  - `<res>` expects any resource location.
  - checks the existance of a `<res>` file.
    - `<res>` is consumed.
    - pushes `true` if file exist, `false` otherwise.
- `<res> touch`
  - `<res>` expects any resource location.
  - checks the existence of a `<res>` file.
    - `<res>` is consumed.
    - pushes `true` if file created or present, `false` if failure.
- `<res> mkdir`
  - `<res>` expects any resource location.
  - checks the existence of a `<res>` directory.
    - `<res>` is consumed.
    - pushes `true` if directory created or present, `false` if failure.
- `<res> rm`
  - `<res>` expects any resource location.
  - removes `<res>` from existence.
    - `<res>` is consumed.
    - pushes `true` if successful, `false` otherwise.
- `<dst> <res> unzip`
  - `<res>` and `<dst>` expects any resource location.
  - unzips `<res>` into `<dst>` directory.
    - both `<res>` and `<dst>` is consumed.
    - pushes `<dir count> <file count>` if successful, `-1 -1` otherwise.
- `<dir> lsf`
  - `<dir>` expects any resource location directory.
  - lists file count in `<dir>`.
    - pushes `<file count>` if successful, `0` otherwise.
- `<idx> <dir> getf`
  - `<dir>` expects any resource location directory.
  - `<idx>` expects a 0-based index within `<dir> lsf` margins.
  - fetches the name of selected file in `<dir>`.
    - pushes string `<name>` if successful, `""` otherwise.
- `<dir> lsd`
  - `<dir>` expects any resource location directory.
  - lists sub-dir count in `<dir>`.
    - pushes `<dir count>` if successful, `0` otherwise.
- `<idx> <dir> getd`
  - `<dir>` expects any resource location directory.
  - `<idx>` expects a 0-based index within `<dir> lsd` margins.
  - fetches the name of selected sub-dir in `<dir>`.
    - pushes string `<name>` if successful, `""` otherwise.
- `<res> <string> concat`
  - `<res>` expects any resource location.
  - `<string>` expects any string.
  - concatenates `<res>` and `<string>` with a directory separator.
    - pushes `<res>`.
- `<string> <string> concat`
  - `<string>` expects any string.
  - concatenates both strings with no separator.
    - pushes concatenated string.
- `<any> tostring`
  - `<any>` expects any standard type.
  - turns any type into a string representation.
    - pushes resulting string if successful, `""` otherwise.
    - if resource location:
      - unless token: leading `/` or `./` is stripped.
      - if token: leading `:` is stripped.
- `<string> token`
  - `<string>` expects any string.
  - turns `<string>` into a token resource location.
    - pushes `:<token>`.
- `<string> absolute`
  - `<string>` expects any string.
  - turns `<string>` into a resource location file relative to `git` root.
    - pushes `/<filepath>`.
- `<string> relative`
  - `<string>` expects any string.
  - turns `<string>` into a resource location file relative to current work dir of script.
    - pushes `./<filepath>`.
- *more to come...*

</details>

<details>
    <summary><code>Language Design</code>: Information about language design</summary>

- A lot of this design follows similarities with Porth (See credits below).

## Branch design
- `while` & `until`:
  - ```
    <prefix> while <condition> do
      <body>
    end
    ```
    - runs `<body>` as long as `<condition>` returns `true`.
      - substitute `while` for `until` and `<body>` runs as long as `<condition>` returns `false`.
- `if` & `unless`:
  - ```
    <condition-0> if
      <body-0>
    end

    <condition-1> if
      <body-1>
    else
      <body-2>
    end
    ```
    - runs `<body-0>` if `<condition-0>` returns `true`.
    - runs `<body-1>` if `<condition-1>` returns `true`, `<body-2>` otherwise.
  - ```
    <condition> unless
      <body-0>
    end

    <condition-1> unless
      <body-1>
    else
      <body-2>
    end
    ```
      - runs `<body-0>` if `<condition-0>` returns `false`.
      - runs `<body-1>` if `<condition-1>` returns `false`, `<body-2>` otherwise.

</details>

# Contributing?
Contributing is not needed, this section is mostly to document the development cycle I follow. But if you do want to contribute, this applies to you as well.

<details>
    <summary><code>Branches</code>: Branches used for development</summary>

### `master`
- Main branch.
- Only updates on new version changes.
- Never pushed directly to.

### `dev`
- Main development branch.
- Merges into master when a new version is ready.
- Never pushed directly to.

### Feature branches
- Main implementation branch.
- Merges into dev when a new feature or fix is ready.
- This is what you should use when contributing.
  - If you've forked, you may still use your `dev` branch in your pull request.

</details>

# Credits
Inspired by these projects:
- [`markut`](https://github.com/tsoding/markut) by [Tsoding/Rexim](https://github.com/tsoding).
- [`porth`](https://gitlab.com/tsoding/porth) by [Tsoding/Rexim](https://gitlab.com/tsoding).
- [premake5](https://premake.github.io/).
- [Hazel `/scripts/` project setup](https://github.com/TheCherno/Hazel/tree/master/scripts) by [TheCherno](https://github.com/TheCherno).
- My extensive use of bash files for development: [ktnlibc/build.sh](https://github.com/Kirdow/ktnlibc/blob/master/build.sh), [cstack/build.sh](https://github.com/Kirdow/cstack/blob/master/build.sh), [pDB/build.sh](https://github.com/Kirdow/pDB/blob/master/run.sh).

# License
`wet` is released under the [MIT License](https://github.com/Ktnuity/wet/blob/master/LICENSE).
