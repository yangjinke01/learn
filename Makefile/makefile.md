# Run and compile your programs more efficiently with this handy automation tool

If you want to run or update a task when certain files are updated, the `make` utility can come in handy. The `make` utility requires a file, `Makefile` (or `makefile`), which defines set of tasks to be executed. You may have used `make` to compile a program from source code. Most open source projects use `make` to compile a final executable binary, which can then be installed using `make install`.

In this article, we'll explore `make` and `Makefile` using basic and advanced examples. Before you start, ensure that `make` is installed in your system.

## Basic examples

Let's start by printing the classic "Hello World" on the terminal. Create a empty directory `myproject` containing a file `Makefile` with this content:

```makefile
say_hello:
        echo "Hello World"
```

Now run the file by typing `make` inside the directory `myproject`. The output will be:

```shell
$ make
echo "Hello World"
Hello World
```

In the example above, `say_hello` behaves like a function name, as in any programming language. This is called the *target*. The *prerequisites* or *dependencies* follow the target. For the sake of simplicity, we have not defined any prerequisites in this example. The command `echo "Hello World"` is called the *recipe*. The *recipe* uses *prerequisites* to make a *target*. The target, prerequisites, and recipes together make a *rule*.

To summarize, below is the syntax of a typical rule:

```
target: prerequisites
<TAB> recipe
```

As an example, a target might be a binary file that depends on prerequisites (source files). On the other hand, a prerequisite can also be a target that depends on other dependencies:

```
final_target: sub_target final_target.c
        Recipe_to_create_final_target

sub_target: sub_target.c
        Recipe_to_create_sub_target
```

It is not necessary for the target to be a file; it could be just a name for the recipe, as in our example. We call these "phony targets."

Going back to the example above, when `make` was executed, the entire command `echo "Hello World"` was displayed, followed by actual command output. We often don't want that. To suppress echoing the actual command, we need to start `echo` with `@`:

```
say_hello:
        @echo "Hello World"
```

Now try to run `make` again. The output should display only this:

```
$ make
Hello World
```

Let's add a few more phony targets: `generate` and `clean` to the `Makefile`:

```makefile
say_hello:
        @echo "Hello World"

generate:
        @echo "Creating empty text files..."
        touch file-{1..10}.txt

clean:
        @echo "Cleaning up..."
        rm *.txt
```

If we try to run `make` after the changes, only the target `say_hello` will be executed. That's because only the first target in the makefile is the default target. Often called the *default goal*, this is the reason you will see `all` as the first target in most projects. It is the responsibility of `all` to call other targets. We can override this behavior using a special phony target called `.DEFAULT_GOAL`.

Let's include that at the beginning of our makefile:

```
.DEFAULT_GOAL := generate
```

This will run the target `generate` as the default:

```
$ make
Creating empty text files...
touch file-{1..10}.txt
```

As the name suggests, the phony target `.DEFAULT_GOAL` can run only one target at a time. This is why most makefiles include `all` as a target that can call as many targets as needed.

Let's include the phony target `all` and remove `.DEFAULT_GOAL`:

```
all: say_hello generate

say_hello:
        @echo "Hello World"

generate:
        @echo "Creating empty text files..."
        touch file-{1..10}.txt

clean:
        @echo "Cleaning up..."
        rm *.txt
```

Before running `make`, let's include another special phony target, `.PHONY`, where we define all the targets that are not files. `make` will run its recipe regardless of whether a file with that name exists or what its last modification time is. Here is the complete makefile:

```
.PHONY: all say_hello generate clean

all: say_hello generate

say_hello:
        @echo "Hello World"

generate:
        @echo "Creating empty text files..."
        touch file-{1..10}.txt

clean:
        @echo "Cleaning up..."
        rm *.txt
```

The `make` should call `say_hello` and `generate`:

```
$ make
Hello World
Creating empty text files...
touch file-{1..10}.txt
```

It is a good practice not to call `clean` in `all` or put it as the first target. `clean` should be called manually when cleaning is needed as a first argument to `make`:

```
$ make clean
Cleaning up...
rm *.txt
```

Now that you have an idea of how a basic makefile works and how to write a simple makefile, let's look at some more advanced examples.

## Advanced examples

### Variables

In the above example, most target and prerequisite values are hard-coded, but in real projects, these are replaced with variables and patterns.

The simplest way to define a variable in a makefile is to use the `=` operator. For example, to assign the command `gcc` to a variable `CC`:

```
CC = gcc
```

This is also called a *recursive expanded variable*, and it is used in a rule as shown below:

```
hello: hello.c
    ${CC} hello.c -o hello
```

As you may have guessed, the recipe expands as below when it is passed to the terminal:

```
gcc hello.c -o hello
```

Both `${CC}` and `$(CC)` are valid references to call `gcc`. But if one tries to reassign a variable to itself, it will cause an infinite loop. Let's verify this:

```
CC = gcc
CC = ${CC}

all:
    @echo ${CC}
```

Running `make` will result in:

```
$ make
Makefile:8: *** Recursive variable 'CC' references itself (eventually).  Stop.
```

To avoid this scenario, we can use the `:=` operator (this is also called the *simply expanded variable*). We should have no problem running the makefile below:

```
CC := gcc
CC := ${CC}

all:
    @echo ${CC}
```

```shell
#$@  表示目标文件
#$^  表示所有的依赖文件
#$<  表示所有的依赖文件中的第一个文件
test : test.go
	go build -o $@ $^
```

```makefile
test:test.o test1.o
    gcc -o $@ $^
    
#"%.o" 把我们需要的所有的 ".o" 文件组合成为一个列表，
#从列表中挨个取出的每一个文件，"%" 表示取出来文件的文件名（不包含后缀），
#然后找到文件中和 "%"名称相同的 ".c" 文件，然后执行下面的命令，直到列表中的文件全部被取出来为止。
%.o:%.c
    gcc -o $@ $^
```

### Functions

* subst
```
$(subst from,to,text)
```

Performs a textual replacement on the text text: each occurrence of from is replaced by to. The result is substituted for the function call. For example,

```
$(subst ee,EE,feet on the street)
```

produces the value ‘fEEt on the strEEt’.

* firstword
```
$(firstword names…)
```

The argument names is regarded as a series of names, separated by whitespace. The value is the first name in the series. The rest of the names are ignored.

For example,

```
$(firstword foo bar)
```

produces the result ‘foo’. Although `$(firstword text)` is the same as `$(word 1,text)`, the `firstword` function is retained for its simplicity.

* shell

  ```
  contents := $(shell cat foo)
  ```

  sets `contents` to the contents of the file foo, with a space (rather than a newline) separating each line.

  ```
  files := $(shell echo *.c)
  ```

  sets `files` to the expansion of ‘*.c’. Unless `make` is using a very strange shell, this has the same result as ‘$(wildcard *.c)’ (as long as at least one ‘.c’ file exists).

* word

```
$(word n,text)
```

Returns the nth word of text. The legitimate values of n start from 1. If n is bigger than the number of words in text, the value is empty. For example,

```
$(word 2, foo bar baz)
```

returns ‘bar’.

* wildcard

```
$(wildcard pattern)
```

The argument pattern is a file name pattern, typically containing wildcard characters (as in shell file name patterns). The result of `wildcard` is a space-separated list of the names of existing files that match the pattern. See [Using Wildcard Characters in File Names](https://www.gnu.org/software/make/manual/html_node/Wildcards.html#Wildcards).

