export function humanSize(size) {
    const suffix = ['B', 'KB', 'MB', 'GB', 'TB']
    let i
    for (i = 0; size > 1024; ++i) {
        size /= 1024
    }
    if (i >= suffix.length) {
        i = suffix.length - 1
    }
    size = Math.round(size * 100) / 100
    return size + ' ' + suffix[i]
}

/**
 *
 * @param {Date} date
 */
export function humanDate(date) {
    let items = [
        date.getFullYear(),
        date.getMonth() + 1,
        date.getDate(),
        date.getHours(),
        date.getMinutes()
    ]
    for (let i = 0; i < items.length; ++i) {
        if (items[i] < 10) {
            items[i] = '0' + items[i]
        }
    }
    return `${items[0]}-${items[1]}-${items[2]} ${items[3]}:${items[4]}`
}
