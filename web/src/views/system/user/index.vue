<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.user']" />
    <OverviewCards ref="overviewRef" class="overview" />
    <a-card class="general-card" :title="$t('menu.system.user')">
      <a-row>
        <a-col :flex="1">
          <a-form :model="query" layout="inline" class="search-form">
            <a-form-item field="keyword" :label="$t('system.common.keyword')">
              <a-input
                v-model="query.keyword"
                :placeholder="$t('system.user.keyword.placeholder')"
                allow-clear
                style="width: 240px"
                @press-enter="search"
              />
            </a-form-item>
            <a-form-item field="roleId" :label="$t('system.user.role')">
              <a-select
                v-model="query.roleId"
                :placeholder="$t('system.user.role.placeholder')"
                allow-clear
                style="width: 160px"
              >
                <a-option v-for="r in roles" :key="r.id" :value="r.id">
                  {{ r.name }}
                </a-option>
              </a-select>
            </a-form-item>
            <a-form-item field="status" :label="$t('system.common.status')">
              <a-radio-group v-model="query.status" type="button">
                <a-radio value="">{{ $t('system.common.all') }}</a-radio>
                <a-radio value="1">{{ $t('system.common.enabled') }}</a-radio>
                <a-radio value="0">{{ $t('system.common.disabled') }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-form>
        </a-col>
        <a-col :flex="'none'">
          <a-space>
            <a-button type="primary" @click="search">
              <template #icon><icon-search /></template>
              {{ $t('system.common.search') }}
            </a-button>
            <a-button @click="reset">
              <template #icon><icon-refresh /></template>
              {{ $t('system.common.reset') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-divider style="margin-top: 4px" />
      <a-row style="margin-bottom: 16px">
        <a-button type="primary" @click="openCreate">
          <template #icon><icon-plus /></template>
          {{ $t('system.user.create') }}
        </a-button>
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :data="rows"
        :pagination="pagination"
        :bordered="false"
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #columns>
          <a-table-column :title="$t('system.user.username')" :width="240">
            <template #cell="{ record }">
              <div class="user-cell">
                <a-avatar :size="36">
                  <img v-if="record.avatar" :src="record.avatar" alt="" />
                  <span v-else>{{ record.username[0]?.toUpperCase() }}</span>
                </a-avatar>
                <div>
                  <div class="user-cell-name">
                    {{ record.nickname || record.username }}
                  </div>
                  <div class="user-cell-sub">@{{ record.username }}</div>
                </div>
              </div>
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.user.role')">
            <template #cell="{ record }">
              <a-tag
                :color="record.role?.code === 'admin' ? 'arcoblue' : 'gray'"
                bordered
              >
                {{ record.role?.name }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.user.email')" data-index="email">
            <template #cell="{ record }">{{ record.email || '-' }}</template>
          </a-table-column>
          <a-table-column :title="$t('system.user.phone')" data-index="phone">
            <template #cell="{ record }">{{ record.phone || '-' }}</template>
          </a-table-column>
          <a-table-column :title="$t('system.common.status')" :width="110">
            <template #cell="{ record }">
              <a-switch
                :model-value="record.status === 1"
                :disabled="record.id === userStore.id"
                size="small"
                @change="(v) => toggleStatus(record, v as boolean)"
              />
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.user.lastLoginAt')">
            <template #cell="{ record }">
              <span v-if="record.lastLoginAt">
                {{ formatTime(record.lastLoginAt) }}
              </span>
              <span v-else class="muted">{{ $t('system.user.never') }}</span>
            </template>
          </a-table-column>
          <a-table-column
            :title="$t('system.common.operations')"
            :width="150"
            fixed="right"
          >
            <template #cell="{ record }">
              <a-button type="text" size="small" @click="openEdit(record)">
                {{ $t('system.common.edit') }}
              </a-button>
              <a-popconfirm
                :content="$t('system.common.deleteConfirm')"
                type="warning"
                @ok="remove(record)"
              >
                <a-button
                  type="text"
                  status="danger"
                  size="small"
                  :disabled="record.id === userStore.id"
                >
                  {{ $t('system.common.delete') }}
                </a-button>
              </a-popconfirm>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="editing ? $t('system.user.edit') : $t('system.user.create')"
      :ok-loading="saving"
      width="560px"
      unmount-on-close
      @before-ok="submit"
    >
      <a-form ref="formRef" :model="form" :rules="rules" auto-label-width>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="username" :label="$t('system.user.username')">
              <a-input v-model="form.username" :disabled="!!editing" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="nickname" :label="$t('system.user.nickname')">
              <a-input v-model="form.nickname" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="password" :label="$t('system.user.password')">
              <a-input-password
                v-model="form.password"
                :placeholder="editing ? $t('system.user.password.editTip') : ''"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="roleId" :label="$t('system.user.role')">
              <a-select
                v-model="form.roleId"
                :disabled="editing?.id === userStore.id"
              >
                <a-option v-for="r in roles" :key="r.id" :value="r.id">
                  {{ r.name }}
                </a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="email" :label="$t('system.user.email')">
              <a-input v-model="form.email" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="phone" :label="$t('system.user.phone')">
              <a-input v-model="form.phone" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="job" :label="$t('system.user.job')">
              <a-input v-model="form.job" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="status" :label="$t('system.common.status')">
              <a-switch
                v-model="form.status"
                :checked-value="1"
                :unchecked-value="0"
                :disabled="editing?.id === userStore.id"
              />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item
              field="introduction"
              :label="$t('system.user.introduction')"
            >
              <a-textarea v-model="form.introduction" :max-length="255" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import type { FieldRule, FormInstance } from '@arco-design/web-vue/es/form';
  import dayjs from 'dayjs';
  import { useUserStore } from '@/store';
  import useLoading from '@/hooks/loading';
  import {
    listUsers,
    listRoles,
    createUser,
    updateUser,
    deleteUser,
    Role,
    SysUser,
    SysUserForm,
  } from '@/api/system';
  import OverviewCards from '../components/overview-cards.vue';

  const { t } = useI18n();
  const userStore = useUserStore();
  const { loading, setLoading } = useLoading();

  const overviewRef = ref<InstanceType<typeof OverviewCards>>();
  const roles = ref<Role[]>([]);
  const rows = ref<SysUser[]>([]);
  const query = reactive({
    keyword: '',
    roleId: undefined as number | undefined,
    status: '',
  });
  const pagination = reactive({
    current: 1,
    pageSize: 10,
    total: 0,
    showTotal: true,
    showPageSize: true,
  });

  const formatTime = (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm');

  const fetchData = async () => {
    setLoading(true);
    try {
      const { data } = await listUsers({
        current: pagination.current,
        pageSize: pagination.pageSize,
        ...query,
      });
      rows.value = data.list;
      pagination.total = data.total;
    } finally {
      setLoading(false);
    }
  };

  const refreshAll = () => {
    fetchData();
    overviewRef.value?.reload();
  };

  const search = () => {
    pagination.current = 1;
    fetchData();
  };
  const reset = () => {
    Object.assign(query, { keyword: '', roleId: undefined, status: '' });
    search();
  };
  const onPageChange = (current: number) => {
    pagination.current = current;
    fetchData();
  };
  const onPageSizeChange = (size: number) => {
    pagination.pageSize = size;
    search();
  };

  const toggleStatus = async (record: SysUser, enabled: boolean) => {
    await updateUser(record.id, {
      ...record,
      status: enabled ? 1 : 0,
    });
    Message.success(t('system.common.success'));
    refreshAll();
  };

  const remove = async (record: SysUser) => {
    await deleteUser(record.id);
    Message.success(t('system.common.success'));
    refreshAll();
  };

  const formRef = ref<FormInstance>();
  const modalVisible = ref(false);
  const saving = ref(false);
  const editing = ref<SysUser | null>(null);
  const emptyForm = (): SysUserForm => ({
    username: '',
    nickname: '',
    password: '',
    email: '',
    phone: '',
    job: '',
    introduction: '',
    status: 1,
    roleId: roles.value.find((r) => r.code === 'user')?.id,
  });
  const form = ref<SysUserForm>(emptyForm());

  const rules = computed<Record<string, FieldRule[]>>(() => ({
    username: [{ required: true, message: t('system.user.username.required') }],
    password: editing.value
      ? [{ minLength: 5, message: t('system.user.password.rule') }]
      : [
          { required: true, message: t('system.user.password.rule') },
          { minLength: 5, message: t('system.user.password.rule') },
        ],
    roleId: [{ required: true, message: t('system.user.role.required') }],
    email: [{ type: 'email', message: t('system.user.email.invalid') }],
  }));

  const openCreate = () => {
    editing.value = null;
    form.value = emptyForm();
    modalVisible.value = true;
  };
  const openEdit = (record: SysUser) => {
    editing.value = record;
    form.value = { ...record, password: '' };
    modalVisible.value = true;
  };

  const submit = async (done: (closed: boolean) => void) => {
    if (await formRef.value?.validate()) {
      done(false);
      return;
    }
    saving.value = true;
    try {
      if (editing.value) await updateUser(editing.value.id, form.value);
      else await createUser(form.value);
      Message.success(t('system.common.success'));
      done(true);
      refreshAll();
    } catch {
      done(false);
    } finally {
      saving.value = false;
    }
  };

  onMounted(async () => {
    const { data } = await listRoles();
    roles.value = data;
    fetchData();
  });
</script>

<script lang="ts">
  export default {
    name: 'SystemUser',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px;
  }

  .overview {
    margin-bottom: 16px;
  }

  .search-form {
    :deep(.arco-form-item) {
      margin-bottom: 16px;
    }
  }

  .user-cell {
    display: flex;
    gap: 12px;
    align-items: center;

    &-name {
      color: var(--color-text-1);
      font-weight: 500;
    }

    &-sub {
      color: var(--color-text-3);
      font-size: 12px;
    }
  }

  .muted {
    color: var(--color-text-4);
  }
</style>
