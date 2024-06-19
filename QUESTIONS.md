# Greystone Candidate Test

**Note: If your IDE does not provide the ability to run specific tests, you can do so with CLI command `go test -run <test-name> ./...` from the project's root directory**

# Question 1: Algorithmic implementation

## Objective

A South African ID number is a 13-digit number which is defined by the following format: YYMMDDSSSSCAZ.

* The first 6 digits (YYMMDD) are based on your date of birth.
* The next 4 digits (SSSS) are used to define your gender.  Females are assigned numbers in the range 0000-4999 and males from 5000-9999.
* The next digit (C) shows if you're an SA citizen status with 0 denoting that you were born a SA citizen and 1 denoting that you're a permanent resident.
* The last digit (Z) is a checksum digit – used to check that the number sequence is accurate using a set formula called the Luhn algorithm.

Complete the function in `questions/question1.go` that parses a given ID number and returns the following:

1. A pointer to an `SAIDDetails` struct and a nil error if the ID is valid
2. An error if the ID is invalid (either through the Luhn checksum or the length)

The Luhn Algorithm follows the below formula:

(Remember that the result of multiplication is called the product and the result of addition is called the sum)

1. Find the sum of all the odd position digits (1st digit, 3rd digit, 5th digit etc) **excluding the checksum digit**
2. Find the sum of all even position digits multiplied by 2 (multiply the digit, not the sum). If the product of that multiplication is greater than 9, then add the 2 digits together (hint, can this be done another way?)
3. Find the sum of steps 1, 2 and the checksum digit. If the total is a multiple of 10, the checksum passes

## Example

Given an ID of `910310 5017 0 84`, we can deduce that this person

1. Was born on the 10th of March 1991
2. Is Male
3. Is a South African Citizen
4. And the ID is valid because the Luhn checksum passes

### Demonstration of Luhn Checksum as per the steps described above

1. 9+0+1+5+1+0 = 16
2. (1x2) + (3x2) + (0x2) + (0x2) + (7x2) + (8x2) = 20 (*remember the rule about products greater than 9*)
3. 16 + 20 + 4 = 40. 40 % 10 = 0

## Instructions

This question has a separate test file. In order to pass, you need to complete the algorithm and make sure your tests pass to.

No `main` package has been included, but you may implement one if you wish. You will only be scored on your algorithm's implementation for this question

# Question 2: Fixing broken code

## Objective

The function `IsFibNumber` is broken. It should verify whether any given integer is a member of the set of Fibonacci numbers.

Fix the function such that both tests pass

## Instructions

This question has a separate test file. In order to pass, you need to complete the algorithm and make sure your tests pass to.

No `main` package has been included, but you may implement one if you wish. You will only be scored on your algorithm's implementation for this question

# Question 3: Concurrency

## Objectives

A common problem in computer science is processing tasks in a concurrent manner.

Question 3 models a Request, some arbitrary unit of work, and a Request Manager.

Complete the Request Manager such that any client can safely call any of the 3 methods without causing a crash, and such that these tasks are truly concurrent.

Note: concurrency means one or more tasks can start, run and complete in overlapping time periods. It doesn't mean they will ever run at the **same instant**, that is parallelism and is not a requirement for this task. This can also be explained as the difference between multi-tasking on a single-core machine (concurrency) and a multi-core machine (parallelism).

## Instructions

There is no test file for this question, you may test this however you choose (I recommend using a main method and really making your solution work hard). Be sure you think extra hard about thread safety.

# Question 4: Systems design

## Objectives

Build a simple RESTFUL http service in Golang with the following (intentionally vague) requirements.

* Login

* Profile Update

* Profile Delete

* Profile Query

The `profile` must contain the following:

* ProfileId

* Name

* Age

* FavoriteColor

* FavoriteOperatingSystem

You may choose whichever supporting technologies/frameworks you wish, this is also a test of your knowledge of technologies outside of Golang. Please try to keep the submission simple. Do not worry if it requires external connections to services. This is an exercise in your ability to plan and layout a small service.

**NB**: It goes without saying that the `login` feature **must** be secure (i.e. passwords must match the profile and should not be stored in an unsafe manner). Implement whichever technique(s) you prefer

We may contact you if we deem it necessary to run the solution to test it thoroughly.
