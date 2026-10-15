<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.role']" />
    <a-card class="general-card" :title="$t('menu.system.role')">
      <template #extra>
        <a-input-search
          v-model="keyword"
          :placeholder="$t('system.role.name')"
          allow-clear
          style="width: 220px"
          @search="fetchData"
          @clear="fetchData"
        />
      </template>
      <a-alert type="info" class="tip">{{ $t('system.role.tip') }}</a-alert>
      <a-spin :loading="loading" style="width: 100%">
        <a-grid
          :cols="{ xs: 1, sm: 2, lg: 3, xxl: 4 }"
          :col-gap="16"
          :row-gap="16"
        >
          <a-grid-item>
            <div class="role-card role-card-create" @click="openCreate">
              <icon-plus :size="24" />
              <span>{{ $t('system.role.create') }}</span>
            </div>
          </a-grid-item>
          <a-grid-item v-for="role in roles" :key="role.id">
            <div class="role-card" :class="{ disabled: role.status !== 1 }">
              <div class="role-card-head">
                <div
                  class="role-card-icon"
                  :class="role.code === 'admin' ? 'primary' : ''"
                >
                  <icon-safe v-if="role.code === 'admin'" />
                  <icon-user v-else />
                </div>
                <div class="role-card-title">
                  <div class="name">
                    {{ role.name }}
                    <a-tag v-if="role.builtin" size="small" color="arcoblue">
                      {{ $t('system.role.builtin') }}
                    </a-tag>
                    <a-tag v-if="role.status !== 1" size="small">
                      {{ $t('system.common.disabled') }}
                    </a-tag>
                  </div>
                  <a-typography-text code class="code">
                    {{ role.code }}
                  </a-typography-text>
                </div>
              </div>
              <div class="role-card-desc">{{ role.description || '-' }}</div>
              <div class="role-card-foot">
                <span class="count">
                  <icon-user-group />
                  {{ role.userCount }} {{ $t('system.role.userCount') }}
                </span>
                <a-space :size="4">
                  <a-button type="text" size="mini" @click="openEdit(role)">
                    <template #icon><icon-edit /></template>
                  </a-button>
                  <a-popconfirm
                    :content="$t('system.common.deleteConfirm')"
                    type="warning"
                    @ok="remove(role)"
                  >
                    <a-button
                      type="text"
                      status="danger"
                      size="mini"
                      :disabled="role.builtin || role.userCount > 0"
                    >
                      <template #icon><icon-delete /></template>
                    </a-button>
                  </a-popconfirm>
                </a-space>
              </div>
            </div>
          </a-grid-item>
        </a-grid>
      </a-spin>
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="editing ? $t('system.role.edit') : $t('system.role.create')"
      unmount-on-close
      @before-ok="submit"
    >
      <a-form ref="formRef" :model="form" :rules="rules" auto-label-width>
        <a-form-item field="code" :label="$t('system.role.code')">
          <a-input v-model="form.code" :disabled="editing?.builtin" />
        </a-form-item>
        <a-form-item field="name" :label="$t('system.role.name')">
          <a-input v-model="form.name" />
        </a-form-item>
        <a-form-item field="description" :label="$t('system.role.description')">
          <a-textarea v-model="form.description" :max-length="255" />
        </a-form-item>
        <a-form-item field="status" :label="$t('system.common.status')">
          <a-switch
            v-model="form.status"
            :checked-value="1"
            :unchecked-value="0"
            :disabled="editing?.builtin"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import type { FormInstance } from '@arco-design/web-vue/es/form';
  import useLoading from '@/hooks/loading';
  import {
    listRoles,
    createRole,
    updateRole,
    deleteRole,
    Role,
  } from '@/api/system';

  const { t } = useI18n();
  const { loading, setLoading } = useLoading();
  const keyword = ref('');
  const roles = ref<Role[]>([]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const { data } = await listRoles({ keyword: keyword.value });
      roles.value = data;
    } finally {
      setLoading(false);
    }
  };

  const remove = async (role: Role) => {
    await deleteRole(role.id);
    Message.success(t('system.common.success'));
    fetchData();
  };

  const formRef = ref<FormInstance>();
  const modalVisible = ref(false);
  const editing = ref<Role | null>(null);
  const form = ref<Partial<Role>>({});
  const rules = computed(() => ({
    code: [
      { required: true, message: t('system.role.code.required') },
      { match: /^[a-z0-9_]+$/, message: t('system.role.code.rule') },
    ],
    name: [{ required: true, message: t('system.role.name.required') }],
  }));

  const openCreate = () => {
    editing.value = null;
    form.value = { code: '', name: '', description: '', status: 1 };
    modalVisible.value = true;
  };
  const openEdit = (role: Role) => {
    editing.value = role;
    form.value = { ...role };
    modalVisible.value = true;
  };

  const submit = async (done: (closed: boolean) => void) => {
    if (await formRef.value?.validate()) {
      done(false);
      return;
    }
    try {
      if (editing.value) await updateRole(editing.value.id, form.value);
      else await createRole(form.value);
      Message.success(t('system.common.success'));
      done(true);
      fetchData();
    } catch {
      done(false);
    }
  };

  onMounted(fetchData);
</script>

<script lang="ts">
  export default {
    name: 'SystemRole',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px;
  }

  .tip {
    margin-bottom: 16px;
  }

  .role-card {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 176px;
    padding: 20px;
    background: var(--color-bg-2);
    border: 1px solid var(--color-neutral-3);
    border-radius: 8px;
    transition: all 0.2s ease;

    &:hover {
      border-color: rgb(var(--primary-4));
      box-shadow: 0 6px 20px rgba(var(--primary-6), 0.1);
    }

    &.disabled {
      opacity: 0.6;
    }

    &-create {
      gap: 8px;
      align-items: center;
      justify-content: center;
      color: var(--color-text-3);
      border-style: dashed;
      cursor: pointer;

      &:hover {
        color: rgb(var(--primary-6));
      }
    }

    &-head {
      display: flex;
      gap: 12px;
      align-items: center;
    }

    &-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      color: var(--color-text-2);
      font-size: 20px;
      background: var(--color-fill-2);
      border-radius: 10px;

      &.primary {
        color: rgb(var(--primary-6));
        background: rgb(var(--primary-1));
      }
    }

    &-title {
      .name {
        display: flex;
        gap: 6px;
        align-items: center;
        color: var(--color-text-1);
        font-weight: 500;
        font-size: 15px;
      }

      .code {
        font-size: 12px;
      }
    }

    &-desc {
      flex: 1;
      margin: 14px 0;
      color: var(--color-text-3);
      font-size: 13px;
      line-height: 1.6;
    }

    &-foot {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-top: 12px;
      border-top: 1px solid var(--color-neutral-3);

      .count {
        display: inline-flex;
        gap: 6px;
        align-items: center;
        color: var(--color-text-2);
        font-size: 13px;
      }
    }
  }
</style>
