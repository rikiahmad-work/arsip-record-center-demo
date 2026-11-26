<template>
  <div>
    <!-- Toolbar Peran -->
    <div class="d-flex flex-wrap gap-2 align-items-end mb-3">
      <div class="rbac-role-select flex-grow-1">
        <label class="form-label mb-1">Pilih Peran</label>
        <select v-model="selectedRoleId" class="form-select" @change="syncPermissions">
          <option :value="null" disabled>— Pilih peran —</option>
          <option v-for="role in roles" :key="role.id" :value="role.id">
            {{ role.name }}{{ role.is_system ? ' (Sistem)' : '' }}
          </option>
        </select>
      </div>

      <button class="btn btn-outline-primary" @click="showAddRole = true">
        <i class="ti ti-plus me-1"></i>Tambah Peran
      </button>
      <button
        class="btn btn-outline-danger"
        :disabled="!selectedRole || selectedRole.is_system"
        :title="selectedRole?.is_system ? 'Peran bawaan sistem tidak dapat dihapus' : 'Hapus peran kustom'"
        @click="showHapusRole = true"
      >
        <i class="ti ti-trash me-1"></i>Hapus Peran
      </button>

      <div class="ms-sm-auto w-100 w-sm-auto d-flex flex-wrap gap-2 mt-2 mt-sm-0">
        <button class="btn btn-light flex-grow-1 flex-sm-grow-0" :disabled="!selectedRole" @click="pilihSemua">
          <i class="ti ti-checks me-1"></i>Pilih Semua
        </button>
        <button class="btn btn-light flex-grow-1 flex-sm-grow-0" :disabled="!selectedRole" @click="kosongkan">
          <i class="ti ti-square me-1"></i>Kosongkan
        </button>
        <button class="btn btn-primary flex-grow-1 flex-sm-grow-0" :disabled="!selectedRole || saving" @click="simpan">
          <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
          <i v-else class="ti ti-device-floppy me-1"></i>Simpan Hak Akses
        </button>
      </div>
    </div>

    <div v-if="!selectedRole" class="alert alert-info">
      Pilih salah satu peran untuk melihat dan mengatur matriks hak akses.
    </div>

    <template v-else>
      <div v-if="selectedRole.is_system" class="alert alert-warning py-2">
        <i class="ti ti-info-circle me-1"></i>
        Peran <strong>{{ selectedRole.name }}</strong> adalah peran bawaan sistem dan tidak dapat dihapus.
        <span v-if="selectedRole.key === 'administrator'">
          Hak akses inti Administrator (<code>dashboard.view</code>, <code>users.view</code>, & <code>roles.manage</code>) selalu aktif dan diproteksi otomatis oleh sistem.
        </span>
      </div>

      <div class="row g-3">
        <div v-for="group in permissionGroups" :key="group.group" class="col-md-6 col-xl-4">
          <div class="card h-100 border shadow-none">
            <div class="card-header bg-light d-flex justify-content-between align-items-center py-2">
              <div class="fw-semibold">
                <i :class="['ti', group.icon, 'me-1']"></i>{{ group.group }}
              </div>
              <button
                type="button"
                class="btn btn-sm btn-link p-0 text-decoration-none"
                @click="toggleGroup(group)"
              >
                {{ isGroupChecked(group) ? 'Kosongkan' : 'Pilih semua' }}
              </button>
            </div>
            <div class="card-body py-2">
              <div v-for="(label, key) in group.items" :key="key" class="form-check mb-1">
                <input
                  :id="`perm-${key}`"
                  v-model="selectedPermissions"
                  class="form-check-input"
                  type="checkbox"
                  :value="key"
                >
                <label class="form-check-label small" :for="`perm-${key}`">
                  {{ label }}
                  <code class="text-muted ms-1">{{ key }}</code>
                  <span
                    v-if="selectedRole.key === 'administrator' && (key === 'dashboard.view' || key === 'users.view' || key === 'roles.manage')"
                    class="badge bg-light-primary text-primary ms-1"
                    title="Hak akses inti ini otomatis selalu aktif untuk Administrator saat disimpan"
                  >
                    Wajib
                  </span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Modal Tambah Peran -->
    <UiAppModal v-if="showAddRole" title="Tambah Peran Baru" size="sm" @close="showAddRole = false">
      <form @submit.prevent="submitAddRole">
        <label class="form-label">Nama Peran</label>
        <input
          v-model="newRoleName"
          class="form-control"
          placeholder="Contoh: Pimpinan"
          maxlength="60"
          required
        >
        <small class="text-muted">Peran baru akan dibuat tanpa hak akses. Atur setelah dibuat.</small>
        <div class="d-flex justify-content-end gap-2 mt-3">
          <button type="button" class="btn btn-secondary" @click="showAddRole = false">Batal</button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
            Simpan
          </button>
        </div>
      </form>
    </UiAppModal>

    <!-- Modal Hapus Peran -->
    <UiAppModal v-if="showHapusRole" title="Hapus Peran" size="sm" @close="showHapusRole = false">
      <p class="mb-2">
        Hapus peran <strong>{{ selectedRole?.name }}</strong>?
      </p>
      <p class="text-muted small mb-0">
        Peran kustom yang dihapus tidak dapat dipulihkan. Pastikan tidak ada pengguna yang masih memakai peran ini.
      </p>
      <div class="d-flex justify-content-end gap-2 mt-3">
        <button type="button" class="btn btn-secondary" @click="showHapusRole = false">Batal</button>
        <button type="button" class="btn btn-danger" :disabled="saving" @click="submitHapusRole">
          <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
          Hapus
        </button>
      </div>
    </UiAppModal>
  </div>
</template>

<script setup lang="ts">
export interface RoleItem {
  id: number
  key: string
  name: string
  is_system: boolean
  users_count: number
  permissions: string[]
}

export interface PermissionGroup {
  group: string
  icon: string
  items: Record<string, string>
}

interface Props {
  roles: RoleItem[]
  permissionGroups: PermissionGroup[]
  saving?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  saving: false,
})

const emit = defineEmits<{
  save: [roleId: number, permissions: string[]]
  addRole: [name: string]
  deleteRole: [roleId: number]
}>()

const selectedRoleId = ref<number | null>(null)
const selectedPermissions = ref<string[]>([])
const showAddRole = ref(false)
const showHapusRole = ref(false)
const newRoleName = ref('')

const selectedRole = computed<RoleItem | null>(() =>
  props.roles.find((r) => r.id === selectedRoleId.value) ?? null
)

// Pilih role pertama secara otomatis saat data tersedia
watch(
  () => props.roles,
  (roles) => {
    if (selectedRoleId.value === null && roles.length > 0) {
      const first = roles[0]
      if (first) {
        selectedRoleId.value = first.id
        syncPermissions()
      }
    } else if (selectedRoleId.value !== null) {
      syncPermissions()
    }
  },
  { immediate: true }
)

function syncPermissions() {
  const role = selectedRole.value
  selectedPermissions.value = role ? [...role.permissions] : []
}

function allKeys(): string[] {
  return props.permissionGroups.flatMap((g) => Object.keys(g.items))
}

function pilihSemua() {
  selectedPermissions.value = allKeys()
}

function kosongkan() {
  selectedPermissions.value = []
}

function isGroupChecked(group: PermissionGroup): boolean {
  const keys = Object.keys(group.items)
  return keys.length > 0 && keys.every((k) => selectedPermissions.value.includes(k))
}

function toggleGroup(group: PermissionGroup) {
  const keys = Object.keys(group.items)
  if (isGroupChecked(group)) {
    selectedPermissions.value = selectedPermissions.value.filter((k) => !keys.includes(k))
  } else {
    const merged = new Set([...selectedPermissions.value, ...keys])
    selectedPermissions.value = Array.from(merged)
  }
}

function simpan() {
  if (!selectedRole.value) return
  emit('save', selectedRole.value.id, [...selectedPermissions.value])
}

function submitAddRole() {
  const name = newRoleName.value.trim()
  if (!name) return
  emit('addRole', name)
  newRoleName.value = ''
  showAddRole.value = false
}

function submitHapusRole() {
  if (!selectedRole.value) return
  emit('deleteRole', selectedRole.value.id)
  showHapusRole.value = false
}
</script>

<style scoped>
.rbac-role-select {
  min-width: 240px;
}
@media (max-width: 575.98px) {
  .rbac-role-select {
    width: 100% !important;
    min-width: 100% !important;
  }
}
</style>
