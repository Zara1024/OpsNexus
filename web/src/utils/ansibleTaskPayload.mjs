export function buildRecurringTogglePayload(row) {
  return {
    id: row.id,
    isRecurring: Number(row.is_recurring),
    cronExpr: row.cron_expr || ''
  }
}
