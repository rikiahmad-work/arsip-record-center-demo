/**
 * Utility Pemformatan Tanggal & Waktu Bahasa Indonesia (id-ID)
 * Format baku kedinasan Republik Indonesia / Kemen PPPA
 */

const MONTHS_INDO = [
  'Januari',
  'Februari',
  'Maret',
  'April',
  'Mei',
  'Juni',
  'Juli',
  'Agustus',
  'September',
  'Oktober',
  'November',
  'Desember',
]

const MONTHS_SHORT_INDO = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'Mei',
  'Jun',
  'Jul',
  'Agu',
  'Sep',
  'Okt',
  'Nov',
  'Des',
]

/**
 * Format tanggal dalam Bahasa Indonesia (contoh: 10 Mei 2026)
 * Aman dari UTC timezone shift untuk format YYYY-MM-DD
 */
export function formatDateIndo(
  val: string | Date | null | undefined,
  options?: { shortMonth?: boolean }
): string {
  if (!val) return '-'

  const months = options?.shortMonth ? MONTHS_SHORT_INDO : MONTHS_INDO

  // Jika berupa string YYYY-MM-DD murni
  if (typeof val === 'string') {
    const trimmed = val.trim()
    const match = trimmed.match(/^(\d{4})-(\d{2})-(\d{2})/)
    if (match) {
      const year = match[1]
      const monthIdx = parseInt(match[2], 10) - 1
      const day = parseInt(match[3], 10)
      if (monthIdx >= 0 && monthIdx < 12) {
        return `${day} ${months[monthIdx]} ${year}`
      }
    }
  }

  const d = val instanceof Date ? val : new Date(val)
  if (isNaN(d.getTime())) return String(val)

  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: options?.shortMonth ? 'short' : 'long',
    year: 'numeric',
  })
}

/**
 * Format tanggal & waktu lengkap dalam Bahasa Indonesia (contoh: 25 November 2026, 15:14 WIB)
 */
export function formatDateTimeIndo(
  val: string | Date | null | undefined,
  options?: { withSeconds?: boolean, shortMonth?: boolean }
): string {
  if (!val) return '-'

  // Jika input hanya YYYY-MM-DD tanpa komponen jam
  if (typeof val === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(val.trim())) {
    return formatDateIndo(val, options)
  }

  const d = val instanceof Date ? val : new Date(val)
  if (isNaN(d.getTime())) return String(val)

  const datePart = d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: options?.shortMonth ? 'short' : 'long',
    year: 'numeric',
  })

  const timeOptions: Intl.DateTimeFormatOptions = {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }
  if (options?.withSeconds) {
    timeOptions.second = '2-digit'
  }

  const timePart = d.toLocaleTimeString('id-ID', timeOptions).replace(/\./g, ':')

  return `${datePart}, ${timePart} WIB`
}

/**
 * Format waktu saja dalam Bahasa Indonesia (contoh: 15:14 WIB)
 */
export function formatTimeIndo(
  val: string | Date | null | undefined,
  options?: { withSeconds?: boolean }
): string {
  if (!val) return '-'
  const d = val instanceof Date ? val : new Date(val)
  if (isNaN(d.getTime())) return String(val)

  const timeOptions: Intl.DateTimeFormatOptions = {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }
  if (options?.withSeconds) {
    timeOptions.second = '2-digit'
  }

  const timePart = d.toLocaleTimeString('id-ID', timeOptions).replace(/\./g, ':')
  return `${timePart} WIB`
}
