export interface Movie {
  id: number
  movie_id: string

  title: string
  description: string

  genre: string[]

  duration: number
  language: string
  subtitle: string

  rating: string

  poster_url: string
  trailer_url: string

  release_date: string

  status: 'NOW_SHOWING' | 'COMING_SOON' | 'ENDED'

  created_at?: string
  updated_at?: string
}