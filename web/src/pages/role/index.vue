<script setup lang="ts">
import type { IUserProfile } from '@/interface'
import { defineRouteMeta, useRouter } from '@fesjs/fes'
import {
  FButton,
  FForm,
  FFormItem,
  FInput,
  FMessage,
  FModal,
  FOption,
  FPagination,
  FRadio,
  FRadioGroup,
  FSelect,
  FTable,
  FTableColumn,
} from '@fesjs/fes-design'
import { LoadingOutlined } from '@fesjs/fes-design/icon'
import { reactive, ref, watch } from 'vue'
import { useRequest } from 'vue-hooks-plus'
import { request } from '@/api'
import { formatTimestamp } from '@/common/utils'
import { LOGIN_TYPE } from '@/enums'

defineRouteMeta({
  name: 'role/list',
  title: 'Role管理',
})

// 1. Define the Form State matching your Go Model
const createFormState = reactive({
  name: '运营经理',
  key: 'operation_manager',
  status: 1, // Default 1: Enabled
  permission_ids: [] as number[], // 使用 number 数组存储选中的权限ID
})

const pageState = reactive({
  current_page: 1,
  page_size: 10,
  total: 0,
})

const defaultValue = {
  id: '',
  name: '',
  bind_type: 0,
}

// 2. 新增：存储所有可分配的权限列表 (用于 Modal 里的 Select)
const allPermissions = ref<any[]>([]) // 存储从后端获取的权限数据，包含 ID, Name, Key 等

// 3. 新增：获取所有权限的请求
const { run: getAllPermissions } = useRequest(
  () => request('/permission/list'), // 假设后端有一个 /permission/all 接口返回所有权限
  {
    manual: false,
    onSuccess: (resData) => {
      // 假设后端直接返回 list 数组
      allPermissions.value = resData.list || resData.data.list || []
    },
  },
)

const state = reactive({
  modal: false,
})

// 4. 修正：Modal 打开时加载权限数据
watch(
  () => state.modal,
  (show) => {
    if (show) {
      getAllPermissions() // Modal 打开时加载权限列表
    }
  },
)

const searchState = reactive(defaultValue)
const router = useRouter()

function reset() {
  router.go(0)
}

function reqUserProfile() {
  if (searchState.id && !+searchState.id) {
    FMessage.error('ID 不合法')
    return Promise.reject(new Error('id invalid'))
  }
  return request('/role/list', {
    id: +searchState.id,
    name: searchState.name,
    bind_type: searchState.bind_type,
    ...pageState,
  })
}

const {
  loading,
  data,
  run: getRoleList,
} = useRequest<{ list: IUserProfile[], total: number }>(
  reqUserProfile,
)

const loadingOnce = ref(loading.value)
const formRef = ref<typeof FForm>()
watch(
  () => loading.value,
  val => (loadingOnce.value = val),
  { once: true },
)

function handleChange(page: number, pageSize: number) {
  pageState.current_page = page
  pageState.page_size = pageSize
  getRoleList()
}

// 5. 修正：handleCreateRole 逻辑 (提交权限 ID)
async function handleCreateRole() {
  // 构造提交的 Body (包含权限 IDs)
  const payload = {
    name: createFormState.name,
    key: createFormState.key,
    status: createFormState.status,
    permission_ids: createFormState.permission_ids, // 🔥 提交选中的 ID 列表
  }

  // 注意：你现在必须保证后端 /role 接口能接收这个新的 payload (包含 permission_ids)
  await request('/role', payload, { method: 'post' })
  FMessage.success('创建成功')
  state.modal = false
  getRoleList()

  // Reset Form
  createFormState.name = ''
  createFormState.key = ''
  createFormState.status = 1
  createFormState.permission_ids = [] // 重置权限列表
}

// 🔥 核心新增：处理权限变更的 PUT 请求
async function handlePermissionChange(roleId: number, newPermIds: number[]) {
  // 1. 构造 PUT 的 Payload
  const payload = {
    permission_ids: newPermIds,
    id: roleId,
  }

  // 2. 调用后端专用 PUT 接口
  // 我们假设后端已经实现了一个 PUT /v1/role/{id}/permissions 接口
  await request(`/role`, payload, { method: 'PUT' })
  FMessage.success('权限分配成功！')

  // 3. 优化：局部刷新
  // 既然更新成功了，我们手动更新前端列表数据，防止全表刷新
  const updatedRole = data.value?.list.find(r => r.id === roleId)
  if (updatedRole) {
    // 注意：因为我们没有获取权限对象的 name/key，所以我们手动用 ID 列表更新当前行的 permissions 属性
    // (这是客户端优化，实际项目中应该让后端返回完整的更新后的 Role 对象)
    updatedRole.permissions = newPermIds.map((id) => {
      // 找到对应的权限对象，保持数据完整性
      const perm = allPermissions.value.find(p => p.id === id)
      return { id, name: perm?.name, key: perm?.key } // 保持 table 结构不崩溃
    })
  }
}
</script>

<template>
  <nav>
    <h1>账号资料</h1>
    <div>
      <FForm
        ref="formRef" :model="data" label-position="right" :span="12" align="flex-start"
        class="user-profile-search-form" @keydown.enter="getRoleList"
      >
        <FFormItem prop="id" label="ID:">
          <FInput v-model="searchState.id" placeholder="请输入ID" @input="pageState.current_page = 1" />
        </FFormItem>
        <FFormItem prop="name" label="姓名:">
          <FInput v-model="searchState.name" placeholder="请输入姓名" @input="pageState.current_page = 1" />
        </FFormItem>
        <FFormItem prop="login_type" label="绑定类型:">
          <FSelect v-model="searchState.bind_type">
            <FOption :value="0">
              全部
            </FOption>
            <FOption
              v-for="(id) in Object.keys(LOGIN_TYPE).filter((k) => isNaN(+(LOGIN_TYPE[k as any])))" :key="id"
              :value="+id"
            >
              {{ LOGIN_TYPE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem style="float: right" label=" ">
          <FButton type="success" @click="state.modal = true">
            创建
          </FButton>
          &nbsp;&nbsp;&nbsp;
          <FButton type="primary" @click="getRoleList">
            查询
          </FButton>
          &nbsp;&nbsp;&nbsp;
          <FButton @click="reset">
            重置
          </FButton>
        </FFormItem>
      </FForm>
    </div>
  </nav>
  <div v-if="loading" class="loading">
    <LoadingOutlined class="icon" />
  </div>
  <FTable
    v-show="!loading" always-scrollbar class="table" :height="10" size="small" row-key="id"
    :data="data?.list ?? []"
  >
    <FTableColumn fixed="left" prop="id" label="Role ID" :min-width="60" />
    <FTableColumn prop="name" label="角色名称" :min-width="150" />
    <FTableColumn label="权限分配/操作" :min-width="350">
      <template #default="{ row }">
        <FSelect
          multiple filterable placeholder="分配权限" :model-value="row.permissions?.map((p: any) => p.id)"
          :options="allPermissions" value-field="id" label-field="name"
          @change="(newIds: number[]) => handlePermissionChange(row.id, newIds)"
        />
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="created_at" label="创建时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.created_at) }}
      </template>
    </FTableColumn>
  </FTable>
  <FPagination
    v-if="!loadingOnce" class="pagination" show-total :total-count="data?.total" show-size-changer
    show-quick-jumper :page-size="pageState.page_size" @change="handleChange"
  />
  <FModal v-model:show="state.modal" title="创建Role" display-directive="show" @ok="handleCreateRole">
    <FForm
      ref="formRef" :model="createFormState" label-position="top" :span="12" align="flex-start"
      class="user-profile-search-form1"
    >
      <FFormItem prop="name" label="角色名称:">
        <FInput v-model="createFormState.name" placeholder="例如：运营经理" />
      </FFormItem>

      <FFormItem prop="key" label="角色标识:">
        <FInput v-model="createFormState.key" placeholder="例如：operation_manager" />
      </FFormItem>

      <FFormItem prop="permission_ids" label="分配权限:">
        <FSelect v-model="createFormState.permission_ids" placeholder="请选择角色权限" multiple filterable>
          <FOption
            v-for="perm in allPermissions" :key="perm.id" :value="perm.id"
            :label="`${perm.name} (${perm.key})`"
          />
        </FSelect>
      </FFormItem>

      <FFormItem prop="status" label="状态:">
        <FRadioGroup v-model="createFormState.status">
          <FRadio :value="1">
            启用
          </FRadio>
          <FRadio :value="0">
            禁用
          </FRadio>
        </FRadioGroup>
      </FFormItem>
    </FForm>
  </FModal>
</template>

<style scoped lang="less">
.loading {
  width: 100%;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;

  .icon {
    margin-top: 20px;
    font-size: 40px;
  }
}

.preview-image {
  width: 30px;
  height: 30px;
  margin-right: 3px;
  display: inline-block;
  overflow: hidden;
  border: 1px solid #333;
  border-radius: 3px;
}

.pagination {
  margin-top: 10px;
  align-self: center;
}

nav {
  margin-bottom: 20px;
}

.table {
  flex: 1;
  display: flex;
  flex-direction: column;

  :global(.table .fes-table-body-wrapper) {
    flex: 1;
    overflow: scroll;
  }
}

.user-profile-search-form {
  :global(& .fes-form-item) {
    // display: inline-block;
    display: inline-flex;
    flex-direction: row;

    &+& {
      margin-right: 24px;
    }
  }

  :global(& .fes-input) {
    width: 150px;
  }

  :global(& .fes-select) {
    width: 100px;
  }
}
</style>
