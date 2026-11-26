// ================================================================
// SiARSIP — Core Role, Session, Menu, & Row Types (Task 4.1)
// ================================================================

export type AdminRoleKey =
  | 'administrator'
  | 'operator'
  | 'viewer'
  | 'approver'
  | 'super-admin'
  | 'admin-arsip'

export interface AdminUserSession {
  id?: number
  username: string
  name: string
  nip: string
  email: string
  roleKey: AdminRoleKey
  roleLabel: string
  unit: string
  jabatan: string
  isSso: boolean
  permissions?: string[]
}

export interface AdminMenuItem {
  key: string
  label: string
  to: string
  icon: string
  group?: string
  roles: AdminRoleKey[]
  permission?: string
}

export interface ArsipRow {
  id: string
  noItemArsip: string
  noFolder: string
  kodeKlasifikasi: string
  uraianInformasiArsip: string
  jenisArsip: string
  kurunWaktu: string
  tingkatPerkembangan: string
  jumlah: string
  keterangan: string
  nomorDefinitif: string | null
  lokasiSimpan: string
  lokasiSimpanId?: number | null
  noBoks: string
  noRak: string
  jangkaSimpanInaktif: string
  jangkaSimpanNasibAkhir: string
  nasibAkhir: string
  skkad: string
  kategoriArsip: string
  penciptaArsip: string
  unitKerjaId?: number | null
  [key: string]: any
}

export interface UserRow {
  id: number
  username: string
  nama: string
  nip: string
  email: string
  unitKerja: string
  jabatan: string
  role: string
  roleLabel: string
  isSso: boolean
  isActive: boolean
  mfaEnabled?: boolean
  lastLoginAt: string | null
  createdAt: string
  [key: string]: any
}

export interface KodeKlasifikasi {
  id: string
  kode: string
  nama: string
  jangkaSimpanInaktif?: string | null
  nasibAkhir?: string | null
  skkad?: string
  totalArsip: number
  statusKey: 'aktif' | 'nonaktif'
  statusLabel: string
}
