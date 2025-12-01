<script lang="ts" setup>
import { useModel, useRouter } from '@fesjs/fes'
import { FButton, FAvatar } from '@fesjs/fes-design'
import { setToken } from '@/common/utils'

const router = useRouter()
async function handleLogout() {
  setToken('')
  router.replace('/login')
}
const initialState = useModel('@@initialState')
function handleClick() {
  router.push('/profile')
}
</script>

<template>
  <div class="right">
    <a class="avatar-container" @click="handleClick">
      <FAvatar size="small" :src="initialState?.image" class="avatar" alt="avatar" />
      <span class="name">
        {{ initialState.nickname || initialState.email || initialState.userId }}
      </span>
    </a>
    <FButton type="link" class="link" @click="handleLogout">
      退出登录
    </FButton>
  </div>
</template>

<style scope lang="less">
.link {
  color: rgb(255, 255, 255);
}
.right {
  text-align: right;
  padding: 0 20px;
}
.avatar-container {
  display: inline-flex;
  flex-direction: column;
  width: 60px;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  .avatar {
    border-radius: 50%;
    display: inline-block;
  }
  .name {
    color: #fff;
    font-size: 12px;
    line-height: 12px;
    margin-top: 4px;
  }
}
</style>
