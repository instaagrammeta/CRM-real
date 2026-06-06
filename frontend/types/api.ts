// =============================================================
// API types — синхронизированы с Go backend (internal/models)
// =============================================================

export interface User {
  id: number
  full_name: string
  age: number
  personal_phones: string[] | null
  work_phones: string[] | null
  login: string
  photo: string
  category: string
  role: 'admin' | 'employee' | string
  telegram_chat_id?: number
  created_at: string
  updated_at: string
}

export interface LoginResponse {
  success: boolean
  role: string
  token: string
  expires: string
  user: Pick<User, 'id' | 'full_name' | 'login' | 'role' | 'category' | 'photo'>
  error?: string
}

export interface Task {
  id: number
  title: string
  description: string
  author_id: number
  executor_id: number
  photo: string
  status: 'new' | 'progress' | 'done' | string
  created_date: string
  updated_at: string
}

export interface Lid {
  id: number
  client_name: string
  phone: string
  topic: string
  comment: string
  source: string
  mortgage: boolean
  box: boolean
  author_id: number
  created_date: string
}

export interface KanbanBoard {
  id: number
  title: string
  color: string
  is_public: boolean
  author_id: number
  is_archived: boolean
  created_date: string
}

export interface KanbanColumn {
  id: number
  board_id: number
  title: string
  color: string
  order_index: number
  is_archived: boolean
  created_date: string
}

export interface KanbanLead {
  id: number
  column_id: number
  board_id: number
  client_name: string
  phone: string
  topic: string
  comment: string
  source: string
  mortgage: boolean
  box: boolean
  author_id: number
  author_name: string
  order_index: number
  created_date: string
  updated_date: string
}

export interface RequestsBoard {
  id: number
  title: string
  color: string
  is_public: boolean
  author_id: number
  created_date: string
}

export interface RequestsColumn {
  id: number
  board_id: number
  title: string
  color: string
  order_index: number
  created_date: string
}

export interface RequestsItem {
  id: number
  column_id: number
  board_id: number
  property_type: string
  address: string
  area: number
  rooms: number
  windows: number
  floor: number
  total_floors: number
  total_price: number
  price_per_m2: number
  phone: string
  client_name: string
  comment: string
  author_id: number
  author_name: string
  executor_id: number
  executor_name: string
  files: any
  order_index: number
  created_date: string
  updated_date: string
}

export interface House {
  id: number
  title: string
  construction_type: string
  district: string
  address: string
  area: number
  rooms: number
  windows: number
  floor: number
  total_floors: number
  price_per_m2: number
  total_price: number
  developer: string
  contact_phone: string
  has_tech_passport: string
  has_renovation_permit: string
  files: any
  author_id: number
  created_date: string
}

export interface Post {
  id: number
  user_id: number
  user_name: string
  title: string
  description: string
  category: string
  content_type: string
  project: string
  media_path: string
  media_type: string
  link: string
  post_date: string
  created_date: string
  likes: number
  comments: number
  shares: number
  views: number
  reach: number
  is_published: boolean
  published_at: string
  updated_date: string
}

export interface Message {
  id: number
  user_id: number
  user_name: string
  message: string
  file_path: string
  file_name: string
  file_type: string
  created_date: string
}

export interface Notification {
  id: number
  user_id: number
  title: string
  body: string
  type: string
  entity_type: string
  entity_id: number
  link: string
  is_read: boolean
  created_at: string
  updated_at: string
}

export interface TelegramSubscriber {
  id: number
  user_id: number
  chat_id: number
  username: string
  first_name: string
  last_name: string
  is_active: boolean
  link_token: string
  created_at: string
}

export interface TelegramLinkResponse {
  success: boolean
  token: string
  link_url: string
  command: string
}

export interface DashboardStats {
  users: number
  tasks: number
  leads: number
  requests: number
  houses: number
  posts: number
  sim_cards: number
}

export interface SimCard {
  id: number
  phone_number: string
  operator: string
  assigned_to: number | null
  phone_id: number | null
  description: string
  status: string
  created_date: string
}

export interface CompanyPhone {
  id: number
  model: string
  phone_id: string
  assigned_to: number | null
  description: string
  status: string
  created_date: string
}

export interface SimTariff {
  id: number
  sim_id: number
  minutes: number
  gb: number
  sms: number
  cost: number
  start_date: string
  end_date: string
  status: string
  created_date: string
}

export interface Bank {
  id: number
  name: string
  slug: string
  logo: string
  description: string
  phone: string
  website: string
  address: string
  order_index: number
  is_active: boolean
  created_date: string
}

export interface InstallmentObject {
  id: number
  name: string
  slug: string
  logo: string
  description: string
  phone: string
  website: string
  address: string
  developer: string
  order_index: number
  is_active: boolean
  created_date: string
}

// =============================================================
// WebSocket events
// =============================================================

export type WsEventType =
  | 'chat:message'
  | 'chat:update'
  | 'chat:delete'
  | 'notification'
  | 'lead:created'
  | 'lead:moved'
  | 'lead:updated'
  | 'task:assigned'
  | 'request:created'
  | 'request:moved'
  | 'tariff:expiring'
  | 'user:online'

export interface WsEvent<T = any> {
  type: WsEventType
  user_id?: number
  data: T
  time: string
}
