## Overview
HourlyPayTracker is a web dashboard app for tracking work hours, calculating income based on hourly rates, and managing payments. 

**Purpose**: Provides freelancers and hourly workers with precise income tracking, payment reconciliation, and financial oversight through an intuitive self-contained web interface.

The frontend is statically embedded directly into the Go backend using `//go:embed static/*`, serving both `/dashboard.html` and `/account.html` from a single binary without external file dependencies. It connects to the embedded API at `http://localhost:8080` for data persistence, using precise Decimal.js arithmetic compliant with ABNT NBR 5891 for financial calculations.

Tested on ARM architecture as well (for fun). If compiled and configured correctly, the complete application can run directly as a process on Android—even through ADB—making the phone act like a server for it. No need for CGO.

The app provides an intuitive interface to log entries, visualize unpaid/paid statuses, and optimize payment matching.

## Key Features
- **Time Entry Logging**: Add/edit entries with clock-in/out datetime pickers, auto-calculating total minutes, hours, and monetary value based on configurable hourly rate.
- **Income Calculation**: Uses high-precision Decimal.js (28 digits, HALF_EVEN rounding) to compute values from hours worked × hourly rate / 60, with cumulative totals.
- **Visual Calendar**: Interactive monthly calendar highlights unpaid (green), paid (purple), and today dates; click days to view/edit entries.
- **Payment Management**: Mark entries as paid individually or in batches; delete selected unpaid/paid entries with confirmation.
- **Smart Received Matching**: Enter received amount to auto-select closest unpaid entries matching the total, shown in a modal with options to mark as paid.
- **Stats Dashboards**: Live stats for total unpaid entries/hours/value; monthly summaries for current/previous months (entries, hours, total revenue).
- **Paginated Tables**: Separate tables for unpaid ("Falta Pagar") and paid ("Pagos") entries, sortable by ID/date, with bulk actions.
- **Responsive Design**: GSAP animations, mobile-optimized (FAB logout, pill buttons), backdrop blur, gradients; works on desktop/tablet/mobile.
- **Modals**: Edit entries, view day-specific logs, received results with diff analysis (success/warn/error tones).
- **Deployment:**: Single Go binary contains both frontend (embedded static files) and backend API. Run go run main.go to serve on :8080.