ago
===

`ago` is _a go_ program.
The name will be changed in future, maybe.

In short, the program chooses useful words in specific documents that need to be understanded well and helps user to understand them well.
The help will be made by listing the each word and its meaning, and user self
test about the word. The self test result will feedback the important word
election.

This program has made for author's english vocabulary memorizing. Because the
description of the program is general, however, the program can be used for
wider general case. For example, terminologies for specific area such as
computer science or mathematics.


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
   future as user want.
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

Examples The attacks have been conducted by an iPhone-hacking wonder kid.  듣기
시작 (Teentimes) 그 공격은 아이폰 해킹 원더 키드에 의해 행해져 왔다.  It seems
as if a series of hacking incidents are plaguing the nation.  듣기 시작
(Teentimes) 일련의 해킹 사고들이 한국을 괴롭히는 것 같다.  Hacking companies to
"help them" in the long run is still unjustified.  듣기 시작 (Teentimes)
"그들을 돕기 위한" 해킹 회사들은 결국 정당화되지 않는다.  2nd Statement: Hacks
can be used as a powerful and effective artistic expression.  듣기 시작
(Teentimes) 2차 진술: 해킹은 강력하고 효과적인 예술적 표현으로 쓰여질 수 있다.
Having said that Nintendo has been hacked and Sony has been hacked twice.  듣기
시작 (Teentimes) 그렇기는 해도 닌텐도도 해킹을 당했고 소니도 두 번이나 해킹을
당했다.  예문 더보기

$ ./ago test
Let the game begin with []

Ready? (Yes/[No])
y
==============================================
Question 1/6:

 [[ math ]]


Press Enter after you remember the meaning of it:

 Reference: 2 Test success/fail: 0/0

----------------------------------------------
The maning of math was:
Meaning
①수학 ②계산 ③수리

Examples
A new application will help those who struggle with math.  듣기 시작
(Kidstimes) 새로운 어플리케이션이 수학으로 고심하는 사람들을 도와줄 것입니다.
A math teacher in the U.S.  recently turned 100 years old.  듣기 시작
(Kindertimes) 최근에 미국의 한 수학 선생님이 100세가 되었습니다.  I think that
teaching children how to do math without calculators is good, too.  듣기 시작
(Kidstimes) 저는 학생들에게 계산기 없이 수학을 하는 방법을 가르치는 것도 좋다고
생각합니다.  If calculators are banned in math class, then students' scores
will only get worse.  듣기 시작 (Kidstimes) 만약 수학시간에 계산기가 금지되면,
학생들의 점수는 더욱 악화 될 것입니다.  What do you think? Should calculators
be banned in math class? 듣기 시작 (Kidstimes) 당신은 어떻게 생각합니까?
수학시간에 계산기 사용을 해야 합니까? 예문 더보기


----------------------------------------------
Were you understanding it well? (Yes/No)
y
your feedback, Yes applied
----------------------------------------------



==============================================
Question 2/6:

 [[ odd ]]


Press Enter after you remember the meaning of it:

 Reference: 1 Test success/fail: 0/0
```


License
=======

GPL v3


Author
======

SeongJae Park (sj38.park@gmail.com)
