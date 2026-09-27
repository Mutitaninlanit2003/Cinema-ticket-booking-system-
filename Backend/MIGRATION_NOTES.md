# MongoDB → PostgreSQL Migration Notes

This project has been migrated from MongoDB (mongo-driver) to PostgreSQL using GORM.

## What changed

- `go.mod`: removed `go.mongodb.org/mongo-driver`; added `gorm.io/gorm`,
  `gorm.io/driver/postgres`, `github.com/google/uuid`, `github.com/lib/pq`.
- `config/database.go`: now opens a `*gorm.DB` PostgreSQL connection via
  `postgres.Open()`, reading `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`,
  `DB_NAME`, `DB_SSLMODE` from `.env`.
- `.env`: MongoDB variables replaced with the PostgreSQL variables above.
- All models (`models/*.go`) now use `ID uint gorm:"primaryKey"` as the
  PostgreSQL primary key. Each model keeps its original external string
  identifier (e.g. `movie_id`, `booking_id`) as a UUID column so existing
  API request/response shapes are unchanged — these UUIDs are auto-generated
  in a `BeforeCreate` hook if not supplied.
- Relationships are expressed as real foreign keys against those UUID
  columns (e.g. `Showtime.MovieID -> Movie.MovieID`), so `Preload`,
  `Where`, and joins work normally.
- New tables `Screen`, `Seat`, `BookingSeat`, and `Payment` were added to
  properly normalize the schema (the original MongoDB schema stored seat
  numbers as a plain string array with no relational integrity). None of
  these have their own controllers/routes yet, since none existed before —
  they exist so the schema is ready for that work. `Booking.SeatNumbers`
  is still populated directly for backward compatibility; `BookingRepository.Create`
  additionally links any seat numbers that match existing `Seat` rows into
  `BookingSeat`, all inside one transaction (rolled back on any failure).
- `repository/*.go` were rewritten to use GORM (`Create`, `Find`, `Where`,
  `Updates`, `Delete`, `Transaction`) instead of the Mongo driver
  (`InsertOne`, `Find`, `UpdateOne`, `DeleteOne`).
- `cmd/main.go` now calls `config.DB.AutoMigrate(...)` for every model
  right after connecting, before routes are registered.
- Controllers, services, and routes keep their original function
  signatures and endpoints — only their internal calls into the
  repository layer changed.
- No authentication/JWT or password hashing existed in the original
  project, so none was added or changed.

## Relationship fix (round 2)

The first pass declared associations on **both** sides of every relationship
— e.g. `Movie.Showtimes []Showtime` (has many) *and* `Showtime.Movie Movie`
(belongs to), both tagged with the same `foreignKey:MovieID;references:MovieID`.
Because `Movie` itself also has a `MovieID` field (its own business key),
GORM's schema parser got ambiguous about which struct actually owns the
foreign key, and tried to add the FK constraint on `movies` while
`showtimes` didn't exist yet — hence `relation "showtimes" does not exist`.

Fix applied: every parent model (`User`, `Movie`, `Theater`, `Screen`,
`Seat`, `Booking`) now declares **no** back-reference slice field at all.
Only the child model declares the association, as a plain "belongs to"
(FK field + struct reference), always pointing child → parent:

- `Showtime` belongs to `Movie`, `Theater`, `Screen`
- `Screen` belongs to `Theater`
- `Seat` belongs to `Screen`
- `Booking` belongs to `User`, `Showtime`
- `BookingSeat` belongs to `Booking`, `Seat`
- `Payment` belongs to `Booking`
- `AuditLog` is independent (no FK constraint at all)

`Showtime.ScreenID` was also changed from `string` to `*string` since the
screen association is optional — a plain `string` defaults to `""`, which
is not SQL `NULL` and would violate the foreign key constraint.

`AutoMigrate` order in `cmd/main.go` already migrates parents before
children (`User, Movie, Theater, Screen, Seat, Showtime, Booking,
BookingSeat, Payment, AuditLog`), which now works because nothing is
migrated out of order or in both directions anymore.

## Numeric primary-key foreign keys (round 3)

All foreign keys now reference the parent's numeric `ID uint` primary key
instead of its business UUID field, per GORM's default (and PostgreSQL's
recommended) belongs-to convention:

- Every child model gained a hidden `*RefID uint` field (e.g.
  `Showtime.MovieRefID`, `Booking.UserRefID`, `BookingSeat.SeatRefID`) —
  this is the actual FK column, indexed and constrained against the
  parent's `ID`. It is tagged `json:"-"` and never exposed to clients.
- The original UUID field (e.g. `Showtime.MovieID`, `Booking.UserID`)
  is kept, with its original `json` tag, but is now a **denormalized
  copy** of the parent's business UUID rather than the relational key
  itself. It's written once at creation time and is what keeps the API
  request/response shape identical to before this refactor.
- Optional associations (`Showtime.Screen`) use pointer types on both
  sides — `ScreenRefID *uint` and `ScreenID *string` — so an absent
  screen is stored as SQL `NULL`, not `0` / `""`.
- Composite unique indexes that used to be built on UUID columns
  (`Seat` per-screen uniqueness, `BookingSeat` per-booking/seat
  uniqueness) now use the numeric `*RefID` columns instead.

Because clients still send/receive the business UUIDs, something has to
translate "UUID from the request" into "numeric ID for the FK column"
before insert. That translation now lives in the repository layer, inside
the same transaction as the insert, so an unresolvable reference rolls
back cleanly instead of leaving a partial row:

- `ShowtimeRepository.Create` looks up `Movie`/`Theater` by their UUID
  (and `Screen`, if provided) and fills in `MovieRefID`/`TheaterRefID`/
  `ScreenRefID` before calling `tx.Create`.
- `BookingRepository.Create` looks up `User`/`Showtime` by their UUID
  and fills in `UserRefID`/`ShowtimeRefID` before calling `tx.Create`;
  the existing best-effort `BookingSeat` linking now also resolves and
  stores `BookingRefID`/`SeatRefID`.

If a referenced UUID doesn't exist, `Create` now returns a clear error
("movie not found for movie_id: ...") instead of silently writing an
invalid row — this is a necessary, intentional behavior change from the
MongoDB version, which had no referential integrity at all.

No changes were needed in `services/` or `controllers/` — they still
bind/pass the same `models.Booking` / `models.Showtime` structs, and the
new `*RefID` fields are `json:"-"` so clients can't (and don't need to)
set them directly.

`Movie`, `Theater`, `User`, and `AuditLog` are unaffected — they were
already parent/independent tables with no foreign keys of their own.

## Before running

1. Create the database: `createdb cinema_booking` (or update `.env`).
2. Update `.env` with your real PostgreSQL credentials.
3. This sandbox had no Go toolchain and no network access, so the code
   could not be compiled or `go mod tidy`'d here. Run, in the project root:

   ```
   go mod tidy
   go run ./cmd
   ```

   `go mod tidy` will fetch the new dependencies and regenerate `go.sum`
   (which was removed since the old one no longer matches `go.mod`).
