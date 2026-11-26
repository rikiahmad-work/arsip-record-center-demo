import type { AdminMenuItem, AdminRoleKey } from '~/types/index'

const ALL_ROLES: AdminRoleKey[] = [
  'administrator',
  'operator',
  'viewer',
  'approver',
  'super-admin',
  'admin-arsip',
]

const ADMIN_ONLY: AdminRoleKey[] = ['administrator', 'super-admin']

export const ADMIN_MENU_ITEMS: AdminMenuItem[] = [
  {
    key: 'dashboard',
    label: 'Dashboard',
    to: '/dashboard',
    icon: 'ti ti-dashboard',
    group: 'Utama',
    roles: ALL_ROLES,
    permission: 'dashboard.view',
  },
  {
    key: 'pengelolaan-arsip',
    label: 'Pengelolaan Arsip',
    to: '/pengelolaan-arsip',
    icon: 'ti ti-file-text',
    group: 'Arsip',
    roles: ALL_ROLES,
    permission: 'archives.view',
  },
  {
    key: 'peminjaman-arsip',
    label: 'Peminjaman Arsip',
    to: '/peminjaman-arsip',
    icon: 'ti ti-archive',
    group: 'Arsip',
    roles: ALL_ROLES,
    permission: 'loans.view',
  },
  {
    key: 'penyusutan-arsip',
    label: 'Penyusutan Arsip',
    to: '/penyusutan-arsip',
    icon: 'ti ti-trash',
    group: 'Arsip',
    roles: ALL_ROLES,
    permission: 'disposal.view',
  },
  {
    key: 'laporan',
    label: 'Laporan Arsip',
    to: '/laporan',
    icon: 'ti ti-chart-bar',
    group: 'Pelaporan & Label',
    roles: ALL_ROLES,
    permission: 'reports.view',
  },
  {
    key: 'label',
    label: 'Cetak Label',
    to: '/label',
    icon: 'ti ti-printer',
    group: 'Pelaporan & Label',
    roles: ALL_ROLES,
    permission: 'reports.export',
  },
  {
    key: 'manajemen-user',
    label: 'Manajemen User',
    to: '/master/user',
    icon: 'ti ti-users',
    group: 'Master Data',
    roles: ADMIN_ONLY,
    permission: 'roles.manage',
  },
  {
    key: 'pengaturan',
    label: 'Pengaturan',
    to: '/master/pengaturan',
    icon: 'ti ti-settings',
    group: 'Master Data',
    roles: ALL_ROLES,
    permission: 'settings.view',
  },
]
