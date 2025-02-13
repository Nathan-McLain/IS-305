import re

# Function to check if a year is a leap year
def is_leap_year(year):
    return (year % 4 == 0 and year % 100 != 0) or (year % 400 == 0)

# Function to validate a date
def is_valid_date(day, month, year):
    # Days in each month (default case)
    month_days = {
        1: 31,  2: 28,  3: 31,  4: 30,  5: 31,  6: 30,
        7: 31,  8: 31,  9: 30, 10: 31, 11: 30, 12: 31
    }
    
    # Adjust for leap year in February
    if month == 2 and is_leap_year(year):
        month_days[2] = 29

    # Validate day range for the given month
    return 1 <= day <= month_days.get(month, 0)

# Regular expression pattern for extracting DD/MM/YYYY
date_pattern = r"\b(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/(1\d{3}|2\d{3})\b"

# Sample text containing dates
text = """
 12/05/2023, 29/02/2024, 31/07/2021.
 31/02/2020, 30/02/2021, 31/04/2021, 29/02/2023.
"""

# Find all potential dates
matches = re.findall(date_pattern, text)

# Process and validate each date
for match in matches:
    day, month, year = map(int, match)  # Convert to integers
    if is_valid_date(day, month, year):
        print(f"Valid date: {day:02d}/{month:02d}/{year}")
    else:
        print(f"Invalid date: {day:02d}/{month:02d}/{year}")
