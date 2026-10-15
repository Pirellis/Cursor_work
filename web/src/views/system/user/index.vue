<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.user']" />
    <a-card class="general-card" :title="$t('menu.system.user')">
      <a-row>
        <a-col :flex="1">
          <a-form :model="query" auto-label-width label-align="left">
            <a-row :gutter="16">
              <a-col :xs="24" :md="8">
                <a-form-item field="username" label="用户">
                  <a-input
                    v-model="query.username"
                    placeholder="用户名 / 昵称"
                    allow-clear
                    @press-enter="search"
                  />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="8">
                <a-form-item field="roleId" label="角色">
                  <a-select
                    v-model="query.roleId"
                    :options="roleOptions"
                    placeholder="全部"
                    allow-clear
                  />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :md="8">
                <a-form-item field="status" label="状态">
                  <a-select
                    v-model="query.status"
                    :options="statusOptions"
                    placeholder="全部"
                    allow-clear
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 32px" direction="vertical" />
        <a-col flex="196px" style="text-align: right">
          <a-space>
            <a-button type="primary" @click="search">
              <template #icon><icon-search /></template>查询
            </a-button>
            <a-button @click="reset">
              <template #icon><icon-refresh /></template>重置
            </a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin-top: 0" />
      <a-row style="margin-bottom: 16px">
        <a-col :span="12">
          <a-button type="primary" @click="openCreate">
            <template #icon><icon-plus /></template>新增用户
          </a-button>
        </a-col>
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :data="list"
        :pagination="pagination"
        :scroll="{ x: 1100 }"
        :bordered="false"
        @page-change="onPageChange"
      >
        <template #columns>
          <a-table-column title="用户" data-index="username" :width="240">
            <template #cell="{ record }">
              <a-space>
                <a-avatar
                  :size="32"
                  :style="{ backgroundColor: avatarColor(record.username) }"
                >
                  <img v-if="record.avatar" :src="record.avatar" alt="" />
                  <template v-else>{{
                    (record.nickname || record.username).slice(0, 1)
                  }}</template>
                </a-avatar>
                <div class="user-cell">
                  <div class="user-cell-name">{{
                    record.nickname || record.username
                  }}</div>
                  <div class="user-cell-sub">@{{ record.username }}</div>
                </div>
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="角色" data-index="role">
            <template #cell="{ record }">
              <a-tag
                :color="record.role.code === 'admin' ? 'arcoblue' : 'gray'"
              >
                {{ record.role.name || '-' }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column title="邮箱" data-index="email">
            <template #cell="{ record }">{{ record.email || '-' }}</template>
          </a-table-column>
          <a-table-column title="手机号" data-index="phone">
            <template #cell="{ record }">{{ record.phone || '-' }}</template>
          </a-table-column>
          <a-table-column title="状态" data-index="status">
            <template #cell="{ record }">
              <a-badge
                :status="record.status === 1 ? 'success' : 'danger'"
                :text="record.status === 1 ? '启用' : '禁用'"
              />
            </template>
          </a-table-column>
          <a-table-column title="最近登录" data-index="lastLoginAt">
            <template #cell="{ record }">{{
              fmt(record.lastLoginAt)
            }}</template>
          </a-table-column>
          <a-table-column title="操作" :width="220" align="center">
            <template #cell="{ record }">
              <a-button type="text" size="small" @click="openEdit(record)">
                编辑
              </a-button>
              <a-button type="text" size="small" @click="openReset(record)">
                重置密码
              </a-button>
              <a-popconfirm
                :content="`确定删除用户 ${record.username} 吗？`"
                type="warning"
                @ok="onDelete(record)"
              >
                <a-button type="text" size="small" status="danger">
                  删除
                </a-button>
              </a-popconfirm>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="formVisible"
      :title="editingId ? '编辑用户' : '新增用户'"
      :ok-loading="saving"
      :on-before-ok="onSubmit"
      unmount-on-close
    >
      <a-form ref="formRef" :model="form" :rules="rules" auto-label-width>
        <a-form-item field="username" label="用户名">
          <a-input
            v-model="form.username"
            :disabled="!!editingId"
            placeholder="登录账号"
          />
        </a-form-item>
        <a-form-item v-if="!editingId" field="password" label="初始密码">
          <a-input-password v-model="form.password" placeholder="至少 5 位" />
        </a-form-item>
        <a-form-item field="nickname" label="昵称">
          <a-input v-model="form.nickname" placeholder="显示名称" />
        </a-form-item>
        <a-form-item field="roleId" label="角色">
          <a-select v-model="form.roleId" :options="roleOptions" />
        </a-form-item>
        <a-form-item field="email" label="邮箱">
          <a-input v-model="form.email" />
        </a-form-item>
        <a-form-item field="phone" label="手机号">
          <a-input v-model="form.phone" />
        </a-form-item>
        <a-form-item field="status" label="状态">
          <a-switch
            v-model="form.status"
            :checked-value="1"
            :unchecked-value="0"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="resetVisible"
      title="重置密码"
      :on-before-ok="onResetSubmit"
      unmount-on-close
    >
      <a-form :model="resetForm" auto-label-width>
        <a-form-item label="用户">{{ resetForm.username }}</a-form-item>
        <a-form-item label="新密码" required>
          <a-input-password
            v-model="resetForm.password"
            placeholder="至少 5 位"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, reactive, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import type { FormInstance } from '@arco-design/web-vue/es/form';
  import dayjs from 'dayjs';
  import useLoading from '@/hooks/loading';
  import {
    queryUsers,
    createUser,
    updateUser,
    deleteUser,
    resetUserPassword,
    getRoleOptions,
    UserRecord,
    UserForm,
    RoleRecord,
  } from '@/api/system';

  const { loading, setLoading } = useLoading(true);
  const list = ref<UserRecord[]>([]);
  const roles = ref<RoleRecord[]>([]);
  const query = reactive<{
    username: string;
    roleId?: number;
    status?: number;
  }>({ username: '' });
  const pagination = reactive({ current: 1, pageSize: 10, total: 0 });

  const roleOptions = computed(() =>
    roles.value.map((r) => ({ label: r.name, value: r.id }))
  );
  const statusOptions = [
    { label: '启用', value: 1 },
    { label: '禁用', value: 0 },
  ];

  const palette = ['#165DFF', '#0FC6C2', '#722ED1', '#F77234', '#00B42A'];
  const avatarColor = (name: string) =>
    palette[
      [...name].reduce((s, c) => s + c.charCodeAt(0), 0) % palette.length
    ];
  const fmt = (v?: string | null) =>
    v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-';

  async function fetchData() {
    setLoading(true);
    try {
      const { data } = await queryUsers({
        current: pagination.current,
        pageSize: pagination.pageSize,
        ...query,
      });
      list.value = data.list;
      pagination.total = data.total;
    } finally {
      setLoading(false);
    }
  }

  function search() {
    pagination.current = 1;
    fetchData();
  }
  function reset() {
    query.username = '';
    query.roleId = undefined;
    query.status = undefined;
    search();
  }
  function onPageChange(current: number) {
    pagination.current = current;
    fetchData();
  }

  const formRef = ref<FormInstance>();
  const formVisible = ref(false);
  const saving = ref(false);
  const editingId = ref<number>();
  const emptyForm = (): UserForm => ({
    username: '',
    password: '',
    nickname: '',
    email: '',
    phone: '',
    roleId: roles.value.find((r) => r.code === 'user')?.id,
    status: 1,
  });
  const form = ref<UserForm>(emptyForm());
  const rules = {
    username: [{ required: true, message: '请输入用户名' }],
    password: [{ required: true, minLength: 5, message: '密码至少 5 位' }],
    roleId: [{ required: true, message: '请选择角色' }],
    email: [{ type: 'email' as const, message: '邮箱格式不正确' }],
  };

  function openCreate() {
    editingId.value = undefined;
    form.value = emptyForm();
    formVisible.value = true;
  }
  function openEdit(r: UserRecord) {
    editingId.value = r.id;
    form.value = {
      username: r.username,
      nickname: r.nickname,
      email: r.email,
      phone: r.phone,
      roleId: r.roleId,
      status: r.status,
    };
    formVisible.value = true;
  }
  async function onSubmit() {
    if (await formRef.value?.validate()) return false;
    saving.value = true;
    try {
      if (editingId.value) {
        await updateUser(editingId.value, form.value);
      } else {
        await createUser(form.value);
      }
      Message.success('保存成功');
      fetchData();
      return true;
    } catch {
      return false;
    } finally {
      saving.value = false;
    }
  }

  async function onDelete(r: UserRecord) {
    await deleteUser(r.id);
    Message.success('已删除');
    fetchData();
  }

  const resetVisible = ref(false);
  const resetForm = reactive({ id: 0, username: '', password: '' });
  function openReset(r: UserRecord) {
    Object.assign(resetForm, { id: r.id, username: r.username, password: '' });
    resetVisible.value = true;
  }
  async function onResetSubmit() {
    if (resetForm.password.length < 5) {
      Message.warning('密码至少 5 位');
      return false;
    }
    try {
      await resetUserPassword(resetForm.id, resetForm.password);
      Message.success('密码已重置');
      return true;
    } catch {
      return false;
    }
  }

  getRoleOptions().then(({ data }) => {
    roles.value = data;
  });
  fetchData();
</script>

<script lang="ts">
  export default {
    name: 'SystemUser',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
  .user-cell {
    line-height: 1.3;
    &-name {
      color: var(--color-text-1);
      font-weight: 500;
    }
    &-sub {
      color: var(--color-text-3);
      font-size: 12px;
    }
  }
</style>
