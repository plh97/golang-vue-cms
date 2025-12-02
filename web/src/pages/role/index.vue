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
  name: 'role/list',
  title: 'Role管理',
})

// 1. Define the Form State matching your Go Model
const createFormState = reactive({
  name: '运营经理',
  key: '',
  status: 1, // Default 1: Enabled
})

const pageState = reactive({
  current_page: 1,
  page_size: 20,
  total: 0,
})

const defaultValue = {
  id: '',
  name: '',
  bind_type: 0,
  wechat: '',
  email: '',
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
    wechat: searchState.wechat,
    email: searchState.email,
    page: {
      current_page: pageState.current_page,
      page_count: 0,
      page_size: pageState.page_size, // 获取全部活动
      total: 0,
    },
  }
  return request('/role/list', data, { method: 'post' })
}

const {
  loading,
  data,
  run: getUserProfileList,
} = useRequest<{ list: IUserProfile[], page: { total: number } }>(
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
  getUserProfileList()
}
function handleAfterEnter(e) {
  console.log('[modal.common] handleAfterEnter, e:', e)
}
function handleAfterLeave(e) {
  console.log('[modal.common] handleAfterLeave, e:', e)
}
// 3. Handle Create Submit
async function handleCreateRole() {
  try {
    await request('/role', createFormState, { method: 'PUT' })
    FMessage.success('创建成功')
    state.modal = false
    getUserProfileList() // You might want to rename this function to getRoleList later

    // Reset Form
    createFormState.name = ''
    createFormState.key = ''
    createFormState.status = 1
  }
  catch (error: any) {
    console.error(error)
    // FMessage is usually handled in request interceptor, but just in case
    // FMessage.error(error.message || '创建失败')
  }
  finally {
    createLoading.value = false
  }
}
</script>

<template>
  <nav>
    <h1>账号资料</h1>
    <div>
      <FForm
        ref="formRef" :model="data" label-position="right" :span="12" align="flex-start"
        class="user-profile-search-form" @keydown.enter="getUserProfileList"
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
              v-for="(id) in Object.keys(LOGIN_TYPE).filter((k) => isNaN(+(LOGIN_TYPE[k as any])))"
              :key="id" :value="+id"
            >
              {{ LOGIN_TYPE[+id] }}
            </FOption>
          </FSelect>
        </FFormItem>

        <FFormItem prop="wxID" label="微信:" @input="pageState.current_page = 1">
          <FInput v-model="searchState.wechat" placeholder="请输入微信号" />
        </FFormItem>

        <FFormItem prop="email" label="邮箱:" @input="pageState.current_page = 1">
          <FInput v-model="searchState.email" placeholder="请输入邮箱" />
        </FFormItem>

        <FFormItem style="float: right" label=" ">
          <FButton type="success" @click="state.modal = true">
            创建
          </FButton>
          &nbsp;&nbsp;&nbsp;
          <FButton type="primary" @click="getUserProfileList">
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
    <FTableColumn fixed="left" prop="id" label="用户ID" :min-width="60" />
    <FTableColumn prop="name" label="用户姓名" />
    <FTableColumn :min-width="50" label="性别">
      <template #default="{ row }">
        {{ GENDER[row.gender] }}
      </template>
    </FTableColumn>
    <FTableColumn prop="wechat" label="绑定微信">
      <template #default="{ row }">
        {{ row.wechat || "-" }}
      </template>
    </FTableColumn>
    <FTableColumn prop="email" label="绑定邮箱">
      <template #default="{ row }">
        {{ row.email || "-" }}
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="create_time" label="创建时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.create_time * 1000) }}
      </template>
    </FTableColumn>
    <FTableColumn :min-width="163" prop="update_time" label="最近更新时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.update_time * 1000) }}
      </template>
    </FTableColumn>
  </FTable>
  <FPagination
    v-if="!loadingOnce"
    class="pagination"
    show-total
    :total-count="data?.page?.total"
    show-size-changer
    show-quick-jumper
    :page-size="pageState.page_size"
    @change="handleChange"
  />
  <FModal
    v-model:show="state.modal"
    title="创建Role"
    display-directive="show"
    @ok="handleCreateRole"
    @after-enter="handleAfterEnter"
    @after-leave="handleAfterLeave"
  >
    <FForm
      ref="formRef" :model="data" label-position="top" :span="12" align="flex-start"
      class="user-profile-search-form1" @keydown.enter="getUserProfileList"
    >
      <FFormItem prop="name" label="角色名称:">
        <FInput
          v-model="createFormState.name"
          placeholder="例如：运营经理"
        />
      </FFormItem>

      <FFormItem prop="key" label="角色标识:">
        <FInput
          v-model="createFormState.key"
          placeholder="例如：operation_manager"
        />
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
