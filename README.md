ago
===

`ago` is _a go_ program.
The name will be changed in future, maybe.

In short, the program chooses useful words in specific documents that need to
be understanded well and helps user to understand them well.
The help will be made by listing the each word and its meaning, and let user to
test whether he is understanding the word well in himself. The self test result
will feedback the important word election.

This program has made for author's english vocabulary memorization training.
Because the description of the program is general, however, the program can be
used for wider general case. For example, terminologies for specific area such
as computer science or mathematics.


Workflow
========

At first, the program receives multiple specific documents.
The documents could be pdf, url, etc in future. For now, however, it receives
text file only.
The program counts frequency of each word in the documents. Words that used
frequently be measured as important.

After that, when user requires, the program shows the each word in importance
order.  For each word, the program gives the user to remember the meaning of
the word.  When user let program knows he has spent enough time to remember,
program shows the meaning of the word and few useful informations about the
word to user.  User can let the program knows whether he remembered the meaning
well or not.
If user says he is remembering the meaning well, importance of the word becomes
lower. In other case, the importance becomes higher.


Usage
=====

Basic usage is similar to other familiar CLI tools like `git`, `perf`, etc.
The usage is as below:
```
ago <command> [arguments ...]
```

In short, ago supports several subcommands with optional arguments for thse subcommands. User can use those subcommands to manage documents and words.

Commands are:
- `add-docs`: Put one or more documents under management of `ago`. Currently,
   only text file is supported. User can specify files they want to add by
   giving path to the files as argument. `ago` automatically extracts words
   inside the documents into its internal storage and manage their importance
   with referenced count and user feedback.
- `ls-docs`: List currently added documents. Output of the command is name of
   document and its id. Because different documents may have same name, a
   document should be identified by the id.
- `rm-docs`: Remove one or more documents from `ago`. User should give id(s) of
   documents they want to remove as argument. Removed documents not be
   displayed via `ls-docs` and words under the documents will be removed from
   `ago`, too.`
- `dic`: Search meaning of a word and its usage example using dictionary.
   Currently, it uses Daum English-Korean Dictionary[1] service. However, as
   the design says, it can be alternated to English dictionary or other in
   future as user want. This command is quite useful for non-English native
   terminal familiar hacker. In the case, he can reduce time consumed by
   dictionary searching because terminal is close and dictionary on browser or
   book shelf is far.
- `test`: Start a self test. Within test, `ago` gives a important word and
   gives time for user to remember the meaning of the word. After that, user
   can signal `ago`. `ago` shows real meaning and usage example of the word to
   user by using the `dic` command internally. After that, user can feedback
   `ago` whether he remembered the meaning well or not. The feedback be used to
   calculate the importance of the word.


Usage Example
=============

Below is an example of `ago` usage. It has lots of assumptions not described
here. However, clever one like you may understand the meaning well ;)

```
$ ./ago ls-docs
$ ./ago add-docs ~/linux/README ~/linux/Documentation/SubmittingPatches
$ ./ago ls-docs
0: README
1: SubmittingPatches
$
$ ./ago rm-docs 1
$ ./ago ls-docs
0: README
$
$ ./ago dic hack
Meaning
①난도질하다 ②늙은 말 ③고용된

Examples
 The attacks have been conducted by an iPhone-hacking wonder kid.
 듣기 시작 (Teentimes) 그 공격은 아이폰 해킹 원더 키드에 의해 행해져 왔다.
 It seems as if a series of hacking incidents are plaguing the nation.
 듣기 시작 (Teentimes) 일련의 해킹 사고들이 한국을 괴롭히는 것 같다.
 Hacking companies to "help them" in the long run is still unjustified.
 듣기 시작 (Teentimes) "그들을 돕기 위한" 해킹 회사들은 결국 정당화되지 않는다.
 2nd Statement: Hacks can be used as a powerful and effective artistic
 expression.
 듣기 시작 (Teentimes) 2차 진술: 해킹은 강력하고 효과적인 예술적 표현으로
 쓰여질 수 있다.
 Having said that Nintendo has been hacked and Sony has been hacked twice.
 듣기 시작 (Teentimes) 그렇기는 해도 닌텐도도 해킹을 당했고 소니도 두 번이나
 해킹을 당했다.
 예문 더보기

$ ./ago test
Let the game begin

Ready? (Yes/[No])
y
==============================================
Question 1/5:

 [[ the ]]


Press Enter after you remember the meaning of it:

 Reference: 185 Test success/fail: 0/0

----------------------------------------------
The maning of the was:
Meaning
①그 ②그럴수록 ③더욱더

Example
The problems do not magically disappear, but you may feel better.
 듣기 시작 (Kidstimes) 그 문제들이 마법처럼 사라지지는 않지만, 너는 기분이
 나아질지도 몰라.
 Every point that you do not challenge makes your opponent's argument stronger.
 듣기 시작 (Kidstimes) 당신이 이의를 제기하지 않는 모든 의견은 당신의 상대편의
 주장을 더욱 강하게 만듭니다.
 A bike can become more so if you attach Fontus to the bike.
 듣기 시작 (Kidstimes) 만약 여러분이 자전거에 Fontus를 부착한다면, 그것은 더욱
 그렇게 될 수 있습니다.
 You can donate money to the NoPhone project on a website.
 듣기 시작 (Kidstimes) 당신은 웹사이트에서 노폰 프로젝트에 돈을 기부할 수
 있습니다.
 Once you do, learning it will become fun instead of work.
 듣기 시작 (Kidstimes) 일단 그렇게 하면, 그것을 배우는 것은 일하는 것이 아닌
 즐거움이 될 것입니다.
 예문 더보기

----------------------------------------------
Were you understanding it well? (Yes/No)
y
your feedback, Yes applied
----------------------------------------------



==============================================
Question 2/5:

 [[ to ]]


Press Enter after you remember the meaning of it:

 Reference: 83 Test success/fail: 0/0
```


License
=======

GPL v3


Author
======

SeongJae Park (sj38.park@gmail.com)
