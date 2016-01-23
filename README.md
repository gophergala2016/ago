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


License
=======

GPL v3


Author
======

SeongJae Park (sj38.park@gmail.com)
