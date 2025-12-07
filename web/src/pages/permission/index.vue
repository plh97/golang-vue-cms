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
import { GENDER, LOGIN_TYPE } from '@/enums'

defineRouteMeta({
  name: 'permission/list',
  title: 'Permission管理',
})

// 1. Define the Form State matching your Go Model
// const createFormState = reactive({
//   name: '运营经理',
//   key: 'operation_manager',
//   status: 1, // Default 1: Enabled
// })
const createFormState = reactive({
  // 基础信息
  name: '查看用户列表',       // 权限名称 (例如：新增用户)
  key: 'sys:user:list',        // 权限标识 (例如：sys:user:add)

  // 核心类型与层级
  type: 2,                     // 权限类型 (默认菜单 2 或 按钮 3)
  parentId: 0,                 // 父级ID (默认为 0，如果是子菜单，则填写父级ID)
  sort: 0,                     // 排序

  // 前端路由 (用于 Type=2)
  path: '/system/user',        // 前端路由地址
  component: 'views/system/user/index', // 前端组件路径

  // 后端鉴权 (用于 Type=3)
  api: '/v1/user/list',        // 后端接口路径
  method: 'POST',              // 请求方法 (通常查询用 GET, 但这里我们暂时用 POST 兼容之前的查询接口)

  // ❌ 注意：Permission 表通常不包含 status 字段，该字段是 Role的属性
  // status: 1, 
})

const pageState = reactive({
  current_page: 1,
  page_size: 15,
  total: 0,
})

const defaultValue = {
  id: '',
  name: '',
  bind_type: 0,
}

const state = reactive({
  modal: false,
})

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
  const data = {
    id: +searchState.id,
    name: searchState.name,
    bind_type: searchState.bind_type,
    current_page: pageState.current_page,
    page_size: pageState.page_size, // 获取全部活动
  }
  return request('/permission/list', data)
}

const {
  loading,
  data,
  run: getPermissionList,
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
  getPermissionList()
}

async function handleCreatePermission() {
  await request('/permission', createFormState, { method: 'post' })
  FMessage.success('创建成功')
  state.modal = false
  getPermissionList() // You might want to rename this function to getRoleList later
  // Reset Form
  createFormState.name = ''
  createFormState.key = ''
  // createFormState.status = 1
}

</script>

<template>
  <nav>
    <h1>账号资料</h1>
    <div>
      <FForm ref="formRef" :model="data" label-position="right" :span="12" align="flex-start"
        class="user-profile-search-form" @keydown.enter="getPermissionList">
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
            <FOption v-for="(id) in Object.keys(LOGIN_TYPE).filter((k) => isNaN(+(LOGIN_TYPE[k as any])))" :key="id"
              :value="+id">
              {{ LOGIN_TYPE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem style="float: right" label=" ">
          <FButton type="success" @click="state.modal = true">
            创建
          </FButton>
          &nbsp;&nbsp;&nbsp;
          <FButton type="primary" @click="getPermissionList">
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
  <FTable v-show="!loading" always-scrollbar class="table" :height="10" size="small" row-key="ID"
    :data="data?.list ?? []">
    <FTableColumn fixed="left" prop="ID" label="ID" :min-width="60" />
    <FTableColumn prop="name" label="权限名称" :min-width="120" />
    <FTableColumn prop="key" label="唯一标识" :min-width="150" />

    <FTableColumn prop="type" label="类型" :min-width="80">
      <template #default="{ row }">
        {{ row.type === 1 ? '目录' : row.type === 2 ? '菜单' : '按钮' }}
      </template>
    </FTableColumn>

    <FTableColumn prop="api" label="API路径/组件" :min-width="200" />
    <FTableColumn prop="method" label="方法" :min-width="80" />

    <FTableColumn :min-width="163" prop="CreatedAt" label="创建时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.CreatedAt) }}
      </template>
    </FTableColumn>
  </FTable>
  <FPagination v-if="!loadingOnce" class="pagination" show-total :total-count="data?.total" show-size-changer
    show-quick-jumper :page-size="pageState.page_size" @change="handleChange" />
  <FModal v-model:show="state.modal" title="创建Permission" display-directive="show" @ok="handleCreatePermission">
    <FForm ref="formRef" :model="createFormState" label-position="top" :span="12" align="flex-start"
      class="permission-create-form">
      <FFormItem prop="type" label="权限类型:">
        <FRadioGroup v-model="createFormState.type">
          <FRadio :value="1">目录</FRadio>
          <FRadio :value="2">菜单</FRadio>
          <FRadio :value="3">按钮</FRadio>
        </FRadioGroup>
      </FFormItem>

      <FFormItem prop="name" label="权限名称:">
        <FInput v-model="createFormState.name" placeholder="例如：删除用户" />
      </FFormItem>
      <FFormItem prop="key" label="唯一标识 (Key):">
        <FInput v-model="createFormState.key" placeholder="例如：sys:user:delete" />
      </FFormItem>

      <template v-if="createFormState.type >= 2">
        <FFormItem prop="path" label="前端路由 Path:">
          <FInput v-model="createFormState.path" placeholder="例如：/system/user" />
        </FFormItem>
        <FFormItem prop="component" label="前端组件路径:">
          <FInput v-model="createFormState.component" placeholder="例如：views/system/user/index" />
        </FFormItem>
      </template>

      <template v-if="createFormState.type === 3">
        <FFormItem prop="api" label="后端接口 API:">
          <FInput v-model="createFormState.api" placeholder="例如：/v1/user/:id" />
        </FFormItem>
        <FFormItem prop="method" label="请求方法:">
          <FSelect v-model="createFormState.method">
            <FOption value="GET">GET</FOption>
            <FOption value="POST">POST</FOption>
            <FOption value="PUT">PUT</FOption>
            <FOption value="DELETE">DELETE</FOption>
          </FSelect>
        </FFormItem>
      </template>

      <FFormItem prop="parentId" label="父级权限 ID:">
        <FInput v-model="createFormState.parentId" type="number" placeholder="0 为顶级" />
      </FFormItem>
      <FFormItem prop="sort" label="排序 (Sort):">
        <FInput v-model="createFormState.sort" type="number" placeholder="数字越小越靠前" />
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
