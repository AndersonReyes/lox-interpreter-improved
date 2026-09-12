# Intro

Source book

craftinginterpreters.com/the-lox-language.html.

This project follows the craftinginterpreters book to implement an interpreter in golang. 
In this project though we will implement a few extra features from other languages:

1. We will introduce the option type natively. no value can be nil. Only an
Option[T] can be nil or a T. To make the syntax digestable we will use the ? operator to define an optional.
recap: An optional can be a Option[T] or T? for short.

2. Implement Scalas Try[T] to be a value that is either a T or 
an Error type (similar to golang error type).

3. Be able to implement interfaces later on 
