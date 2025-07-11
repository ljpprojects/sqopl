# SQOPL

## The "Hello, World 🔥" Example

This example demonstrates how to print a unicode character in SQOPL.
Since SQOPL only works with individual bytes, and it uses Go for the interpreter, unicode strinsg aren't impossible.

You will need to push the neccessary bytes indivudually. You can find which bytes to push by running

```bash
echo -n "[UNICODE CHARACTER HERE]" | xxd
```

You must convert the hexadecimal digits to base-10. For the fire emoji, we can use this command:


```bash
echo -n "🔥" | xxd -p | fold -w 2 | sed 's/^/0x/' | xargs printf "%d\n"
```

It outputs every byte outputted by `echo`, which we convert to base-10 and print on multiple lines. It should output:

```
240
159
148
165
```

Doing this for the skull emoji (💀), we get:

```
240
159
146
128
```

What a nifty one-liner!
