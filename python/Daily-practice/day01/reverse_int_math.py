"""
Reverse an integer using math only (no string conversion).

Problem
-------
Given a positive integer n, return its digits reversed.
Example: 12345 -> 54321

Pseudocode
----------
1. Set reversed = 0
2. While n > 0:
       last_digit = n % 10          # peel off rightmost digit
       reversed = reversed * 10 + last_digit
       n = n // 10                  # drop last digit
3. Return reversed

Concepts used
-------------
- Function: named block of reusable logic (`reverse_integer`)
- Loop: `while` repeats until n becomes 0
- Variables: `reversed`, `last_digit`, `n` hold state that changes each iteration
- Operators: `%` (modulo), `//` (integer division), `*` and `+`

Time complexity: O(d) where d = number of digits in n
Space complexity: O(1) — only a few integer variables, no extra data structures
"""


def reverse_integer(n: int) -> int:
    """
    Reverse the digits of a positive integer without converting to a string.

    Args:
        n: A non-negative integer (0 or greater).

    Returns:
        The integer formed by reversing n's digits.

    Examples:
        >>> reverse_integer(12345)
        54321
        >>> reverse_integer(0)
        0
        >>> reverse_integer(5)
        5
    """
    reversed_num = 0

    while n > 0:
        last_digit = n % 10
        reversed_num = reversed_num * 10 + last_digit
        n = n // 10

    return reversed_num


def main() -> None:
    """Run a quick manual demo when this file is executed directly."""
    input_value = 12345
    result = reverse_integer(input_value)
    print(f"The reversed integer of {input_value} is {result}")


if __name__ == "__main__":
    main()
