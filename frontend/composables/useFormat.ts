// useFormat — formatting helpers (даты, валюта).

export const useFormat = () => {
  const formatDate = (s?: string, withTime = false) => {
    if (!s) return ''
    const d = new Date(s)
    if (Number.isNaN(d.getTime())) return s
    const date = d.toLocaleDateString('ru-RU')
    if (!withTime) return date
    const time = d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
    return `${date} ${time}`
  }

  const formatCurrency = (n?: number, currency = '₽') => {
    if (n == null) return ''
    return new Intl.NumberFormat('ru-RU', { maximumFractionDigits: 2 }).format(n) + ' ' + currency
  }

  const formatNumber = (n?: number) => {
    if (n == null) return ''
    return new Intl.NumberFormat('ru-RU').format(n)
  }

  const initials = (name?: string) => {
    if (!name) return '?'
    return name
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map(s => s[0]?.toUpperCase() || '')
      .join('')
  }

  const timeAgo = (s?: string): string => {
    if (!s) return ''
    const d = new Date(s)
    const diff = Math.floor((Date.now() - d.getTime()) / 1000)
    if (diff < 60) return 'ҳозир'
    if (diff < 3600) return Math.floor(diff / 60) + ' дақ. пеш'
    if (diff < 86400) return Math.floor(diff / 3600) + ' соат пеш'
    if (diff < 30 * 86400) return Math.floor(diff / 86400) + ' рӯз пеш'
    return d.toLocaleDateString('ru-RU')
  }

  return { formatDate, formatCurrency, formatNumber, initials, timeAgo }
}
