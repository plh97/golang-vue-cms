<script setup lang="ts">
import type { IUserProfile } from '@/interface'
import { defineRouteMeta, useRouter } from '@fesjs/fes'
import {
  FButton,
  FForm,
  FFormItem,
  FInput,
  FMessage,
  FOption,
  FPagination,
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
  name: 'account/list',
  title: '账号管理',
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
  return request('/user/list', data, { method: 'post' })
}

const {
  loading,
  data,
  run: getUserProfileList,
} = useRequest<{ list: IUserProfile[], page: { total: number } }>(
  reqUserProfile,
)

const {
  data: allRoles,
  // run: getRoleList,
} = useRequest(
  async () => {
    const res = await request('/role/list')
    return res.list
  },
)

// 🔥 核心新增：处理用户角色变更的 API
// 假设后端有一个 PUT /v1/user/{id}/roles 接口来更新用户的角色列表
async function handleUserRoleChange(userId: string, newRoleIds: number[]) {
  // 构造 PUT 的 Payload
  const payload = {
    role_ids: newRoleIds, // 提交选中的角色 ID 列表
    user_id: userId,
  }

  // 调用后端 API 更新用户的角色列表
  await request(`/user`, payload, { method: 'PUT' })
  FMessage.success('用户角色更新成功！')

  // 刷新列表
  getUserProfileList()
}

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
              v-for="(id) in Object.keys(LOGIN_TYPE).filter((k) => isNaN(+(LOGIN_TYPE[k as any])))" :key="id"
              :value="+id"
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
    <FTableColumn label="角色分配" :min-width="250">
      <template #default="{ row }">
        <FSelect
          multiple filterable placeholder="分配用户角色" :model-value="row.roles?.map(r => r.ID)"
          :options="allRoles ?? []" value-field="ID" label-field="name"
          @change="(newIds: number[]) => handleUserRoleChange(row.user_id, newIds)"
        />
      </template>
    </FTableColumn>

    <FTableColumn :min-width="163" prop="create_time" label="创建时间">
      <template #default="{ row }">
        {{ formatTimestamp(row.created_at) }}
      </template>
    </FTableColumn>
  </FTable>

  <!-- <FTable
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
  </FTable> -->
  <FPagination
    v-if="!loadingOnce" class="pagination" show-total :total-count="data?.page?.total" show-size-changer
    show-quick-jumper :page-size="pageState.page_size" @change="handleChange"
  />
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
