"""
pytest tests for reverse_integer.

Run from this folder:
    pytest -v
    pytest -v test_reverse_int_math.py
"""

from reverse_int_math import reverse_integer


# Step 1: one simple test (mirrors your first Go test in reverse_int_math_test.go)
def test_reverse_integer_basic():
    result = reverse_integer(12345)
    expected = 54321
    assert result == expected


# Step 2 (next lesson): uncomment and run after we cover table-driven tests in Go.
#
# import pytest
#
# @pytest.mark.parametrize(
#     "input_value, expected",
#     [
#         (12345, 54321),
#         (5, 5),
#         (0, 0),
#         (1000, 1),
#         (123456789, 987654321),
#         (1200, 21),
#     ],
#     ids=[
#         "positive_number",
#         "single_digit",
#         "zero",
#         "number_with_zeros",
#         "large_number",
#         "number_ending_in_zero",
#     ],
# )
# def test_reverse_integer_cases(input_value, expected):
#     assert reverse_integer(input_value) == expected
