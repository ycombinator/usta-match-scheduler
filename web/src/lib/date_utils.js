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
    const eventDate = new Date(event.date)
    const eventYear = eventDate.getFullYear()
    const eventMonth = eventDate.getMonth()

    return (year == eventYear && month == eventMonth)
}

export function previousMonth(month) {
    return month > 0 ? month - 1 : 11
}

export function nextMonth(month) {
    return month < 11 ? month + 1 : 0
}

export function isEventInDay(year, month, day, event) {
    const eventDate = new Date(event.date)
    return (isEventInMonth(year, month, event) && day == eventDate.getDate())
}

export function doesEventStartInDay (year, month, day, event) {
    return true
    const startYear = event.date.getFullYear()
    const startMonth = event.date.getMonth()
    const startDay = event.date.getDate()

    return (year == startYear && month == startMonth && day == startDay)
}

export function doesEventEndInDay (year, month, day, event) {
    return true
    const endYear = event.end.getFullYear()
    const endMonth = event.end.getMonth()
    const endDay = event.end.getDate()

    return (year == endYear && month == endMonth && day == endDay)
}

export function getPreviousYearMonth(year, month) {
    const prevMonth = getPreviousMonth(month)
    const prevYear = getPreviousMonthYear(year, month)
    return { prevYear, prevMonth }
}

export function getNextYearMonth(year, month) {
    const nextMonth = getNextMonth(month)
    const nextYear = getNextMonthYear(year, month)
    return { nextYear, nextMonth }
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
    const dt = new Date(date)
    let hours = dt.getHours()
    if (hours > 12) {
        hours = hours % 12
    }
    const hoursStr = hours.toString()
    const minutesStr = dt.getMinutes().toString().padEnd(2,"0")

    if (dt.getMinutes() == 0) {
        return hoursStr
    }

    return hoursStr + ":" + minutesStr
}

export function getMonthName(year, month) {
    const date = new Date(year, month, 1)
    return date.toLocaleString('default', { month: 'long' });
}

function getPreviousMonth(month) {
    if (month == 0) {
        return 11
    }

    return month - 1
}

function getPreviousMonthYear(year, month) {
    if (month == 0) {
        return year - 1
    }

    return year
}

function getNextMonth(month) {
    return (month + 1) % 12
}

function getNextMonthYear(year, month) {
    if (month == 11) {
        return year + 1
    }

    return year
}