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
Since this is a clean slate project, nothing has been done so far. But I'm planning to have this complete feature set and more once the project is finished.

## Tokens vs Files
When a task downloads to a file, it would need an absolute file in relation to your `.git` directory location. However, in order to allow for temporary paths, you will be able to use tokens. To create a token file, you write `:<token>` instead of `<file>`. The prefix of `:` indicates a token being used. On the back-end, this would use a `~/.wet/` for storage, or `./.wet` if manually created next to `./.git`.

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
- `<value> @<name> store`: store `<value>` as `@<name>`.
  - consumes `<value>` and `@<name>`.
  - `@<name>` is a constant non-data-type value.
- `@<name> load`: load from `@<name>`.
  - consumes `@<name>`
  - `@<name>` is a constant non-data-type value.
  - pushes `<value> true` if successful and type matches, `false` otherwise.
- `@<name>` is local-scope by default.
  - prefix `<name>` with `.` (`@.<name>`) for global scope.

## Commands
- `<url> <dst> download`
  - `<url>` expects a `string` containing a remote URL.
  - `<dst>` expects any resource location.
  - downloads the file/resource at `<url>` into `<dst>`
    - both `<url>` and `<dst>` is consumed.
    - pushes `true` if successful, `false` if unsuccessful.
- `<src> <dst> move`
  - `<src>` and `<dst>` both expects any resource location.
  - moves `<src>` to `<dst>`
    - both `<src>` and `<dst>` is consumed.
    - pushes `true` if successful, `false` if unsuccessful.
- `<src> <dst> copy`
  - `<src>` and `<dst>` both expects any resource location.
  - copies `<src>` to `<dst>`
    - both `<src>` and `<dst>` is consumed.
    - pushes `true` if successful, `true false` if `<dst>` already exist, `false false` if unsuccessful.
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
- `<res> rm`
  - `<res>` expects any resource location.
  - removes `<res>` from existence.
    - `<res>` is consumed.
    - pushes `true` if successful, `false` otherwise.
- `<dst> <res> unzip`
  - `<res>` and `<dst>` expects any resource location.
  - unzips `<res>` into `<dst>` directory.
    - both `<res>` and `<dst>` is consumed.
    - pushes `<dir count> <file count>` if successful, `-1` otherwise.
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
    - pushes `<res> true` resource location if successful, `false` otherwise.
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
    - pushes `:<token> true` if successful, `false` otherwise.
- `<string> absolute`
  - `<string>` expects any string.
  - turns `<string>` into a resource location file relative to `git` root.
    - pushes `/<filepath> true` if successful, `false` otherwise.
- `<string> relative`
  - `<string>` expects any string.
  - turns `<string>` into a resource location file relative to current work dir of script.
    - pushes `./<filepath> true` if successful, `false` otherwise.
- *more to come...*

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

# Credits
Inspired by these projects:
- [`markut`](https://github.com/tsoding/markut) by [Tsoding/Rexim](https://github.com/tsoding).
- [`porth`](https://gitlab.com/tsoding/porth) by [Tsoding/Rexim](https://gitlab.com/tsoding).
- [premake5](https://premake.github.io/).
- [Hazel `/scripts/` project setup](https://github.com/TheCherno/Hazel/tree/master/scripts) by [TheCherno](https://github.com/TheCherno).
- My extensive use of bash files for development: [ktnlibc/build.sh](https://github.com/Kirdow/ktnlibc/blob/master/build.sh), [cstack/build.sh](https://github.com/Kirdow/cstack/blob/master/build.sh), [pDB/build.sh](https://github.com/Kirdow/pDB/blob/master/run.sh).

# License
`wet` is released under the [MIT License](https://github.com/Ktnuity/wet/blob/master/LICENSE)
