export function daysInMonth(year, month) {
    switch (month) {
        case 0: // January
        case 2: // March
        case 4: // May
        case 6: // July
        case 7: // August
        case 9: // October
        case 11: // December
            return 31;
        case 3: // April
        case 5: // June
        case 8: // September
        case 10: // November
            return 30;
        case 1: // February
            return (year % 4 === 0 && (year % 100 !== 0 || year % 400 === 0)) ? 29 : 28;
        default:
            return -1; // Invalid month
    }
}

export function weeksInMonth(year, month) {
    return Math.ceil(daysInMonth(year, month) / 7)
}

export function isEventInMonth(year, month, event) {
    const startYear = event.start.getFullYear()
    const startMonth = event.start.getMonth()

    const endYear = event.end.getFullYear()
    const endMonth = event.end.getMonth()

    return (year == startYear && month == startMonth)
        || (year == endYear && month == endMonth)
}

export function isEventInDay (year, month, day, event) {
    const startYear = event.start.getFullYear()
    const startMonth = event.start.getMonth()
    const startDay = event.start.getDate()

    const endYear = event.end.getFullYear()
    const endMonth = event.end.getMonth()
    const endDay = event.end.getDate()

    return (year == startYear && month == startMonth && day == startDay)
        || (year == endYear && month == endMonth && day == endDay)
}

export function getPreviousMonth(month) {
    if (month == 0) {
        return 11
    }

    return month - 1
}

export function getPreviousMonthYear(year, month) {
    if (month == 0) {
        return year - 1
    }

    return year
}
export function getNextMonth(month) {
    return (month + 1) % 12
}

export function getNextMonthYear(year, month) {
    if (month == 11) {
        return year + 1
    }

    return year
}

export function monthDaysInFirstWeek(firstDayOfMonthWeekday) {
    return 7 - firstDayOfMonthWeekday
}

export function isSameDay(date1, date2) {
    return isSameMonth(date1, date2)
        && date1.getDate() == date2.getDate()
}

export function isSameMonth(date1, date2) {
    return date1.getFullYear() == date2.getFullYear()
        && date1.getMonth() == date2.getMonth()
}

export function isWeekendDay(date) {
    return (date.getDay() == 0) || (date.getDay() == 6)
}

export function getPaddedTime(date) {
    let hours = date.getHours()
    if (hours > 12) {
        hours = hours % 12
    }
    const hoursStr = hours.toString()
    const minutesStr = date.getMinutes().toString().padEnd(2,"0")

    if (date.getMinutes() == 0) {
        return hoursStr
    }

    return hoursStr + ":" + minutesStr
}

export function getMonthName(year, month) {
    const date = new Date(year, month, 1)
    return date.toLocaleString('default', { month: 'long' });
}