<script setup>
import { defineRouteMeta, useRouter } from '@fesjs/fes'
import { FButton, FForm, FFormItem, FInput, FSpace } from '@fesjs/fes-design'
import { computed, reactive, ref } from 'vue'
import { request } from '@/api'
import { useValidator } from '@/common/hooks'
import { setToken } from '@/common/utils'

defineRouteMeta({
  name: 'login',
  title: '登录',
  layout: {
    navigation: null,
  },
})

const formRef = ref(null)

const state = reactive({
  submitLoading: false,
  submitText: '登录',
  username: 'admin@gmail.com',
  password: '123456',
})
const validator = useValidator(state)
const rules = computed(() => {
  return {
    username: [
      {
        required: true,
        type: 'string',
        message: '用户名不能为空',
        validator,
      },
    ],
    password: [
      {
        required: true,
        type: 'string',
        message: '密码不能为空',
        validator,
      },
    ],
  }
})
const router = useRouter()

async function submitHandler() {
  await formRef.value.validate()
  state.submitLoading = true
  state.submitText = '校验中'
  request('/login', {
    email: state.username,
    password: state.password,
  }, {
    method: 'POST',
  }).then(async (res) => {
    setToken(res.accessToken)
    // window.location.replace('/user')
    router.replace('/user')
  }).catch(() => {
    state.submitLoading = false
    state.submitText = '提交'
  })
}
function goRegister() {
  request('/register', {
    email: state.username,
    password: state.password,
  }, {
    method: 'POST',
  })
}
</script>

<template>
  <div class="container" id="login">
    <FForm ref="formRef" class="form" label-width="140px" label-position="top" :model="state" :rules="rules">
      <FFormItem prop="username" label="用户名">
        <FInput v-model="state.username" class="input" placeholder="请输入用户名" />
      </FFormItem>
      <div style="height: 20px;"></div>
      <FFormItem prop="password" label="密码">
        <FInput v-model="state.password" class="input" placeholder="请输入密码" type="password" />
      </FFormItem>
      <FFormItem label=" ">
        <FSpace classname="submit-button-space">
          <FButton size="large" class="button" type="primary" :loading="state.submitLoading" @click="submitHandler">
            {{ state.submitText }}
          </FButton>
          <FButton size="large" class="button" type="info" @click="goRegister">
            注册
          </FButton>
        </FSpace>
      </FFormItem>
    </FForm>
  </div>
</template>

<style lang="less" scoped>
.title {
  text-align: center;
  color: #000;
  font-size: 40px;
  font-weight: 500;
  margin: 0px;
  font-size: 40px;
  font-weight: 500;
  line-height: normal;
  margin-bottom: 80px;
}

.container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;

  .logo {
    position: absolute;
    top: 34px;
    left: 28px;
    font-size: 32px;
    color: #000;
    font-size: 32px;
    font-style: normal;
    font-weight: 800;
    line-height: normal;
    display: flex;
    align-items: center;

    img {
      width: 41px;
      height: 41px;
      margin-right: 12px;
    }
  }

  .form {
    width: 508px;
    background-color: #fff;
    box-sizing: border-box;
    padding: 30px 42px;
    border-radius: 10px;
    box-shadow: 0 4px 64px 0 rgba(0, 0, 0, 0.05);
    border: 0.5px solid #878787;

    .button {
      // width: 423px;
      width: calc(410px / 2);
      height: 57px;
      // margin-top: 20px;
    }

    :global(#login .fes-input-inner) {
      height: 59px;
      font-size: 14px;
      padding: 0 27px;
    }

    :global(#login .fes-form .fes-form-item .fes-form-item-label.is-required::before) {
      display: none;
    }

    :global(#login .fes-form-item-error) {
      font-size: 14px;
    }

    :global(#login .fes-form .fes-form-item-top .fes-form-item-label) {
      color: #000;
      font-size: 16px;
      font-style: normal;
      font-weight: 400;
      line-height: normal;
    }
  }
}
</style>
