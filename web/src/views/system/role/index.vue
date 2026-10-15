<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.role']" />
    <a-card class="general-card" :title="$t('menu.system.role')">
      <a-row style="margin-bottom: 16px" justify="space-between">
        <a-button type="primary" @click="openCreate">
          <template #icon><icon-plus /></template>新增角色
        </a-button>
        <a-input-search
          v-model="keyword"
          placeholder="搜索角色名称 / 标识"
          style="width: 280px"
          allow-clear
          @search="search"
          @press-enter="search"
          @clear="search"
        />
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :data="list"
        :pagination="pagination"
        :bordered="false"
        @page-change="onPageChange"
      >
        <template #columns>
          <a-table-column title="角色名称" data-index="name">
            <template #cell="{ record }">
              <a-space>
                <icon-safe
                  v-if="record.code === 'admin'"
                  :style="{ color: 'rgb(var(--arcoblue-6))' }"
                />
                <icon-user v-else :style="{ color: 'var(--color-text-3)' }" />
                <span>{{ record.name }}</span>
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="角色标识" data-index="code">
            <template #cell="{ record }">
              <a-tag>{{ record.code }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column
            title="描述"
            data-index="description"
            ellipsis
            tooltip
          >
            <template #cell="{ record }">{{
              record.description || '-'
            }}</template>
          </a-table-column>
          <a-table-column title="用户数" data-index="userCount" :width="100" />
          <a-table-column title="状态" data-index="status" :width="100">
            <template #cell="{ record }">
              <a-badge
                :status="record.status === 1 ? 'success' : 'danger'"
                :text="record.status === 1 ? '启用' : '禁用'"
              />
            </template>
          </a-table-column>
          <a-table-column title="创建时间" data-index="createdAt">
            <template #cell="{ record }">
              {{ dayjs(record.createdAt).format('YYYY-MM-DD HH:mm') }}
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="160" align="center">
            <template #cell="{ record }">
              <a-button type="text" size="small" @click="openEdit(record)">
                编辑
              </a-button>
              <a-popconfirm
                :content="`确定删除角色 ${record.name} 吗？`"
                type="warning"
                @ok="onDelete(record)"
              >
                <a-button
                  type="text"
                  size="small"
                  status="danger"
                  :disabled="record.code === 'admin'"
                >
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
      :title="editingId ? '编辑角色' : '新增角色'"
      :on-before-ok="onSubmit"
      unmount-on-close
    >
      <a-form ref="formRef" :model="form" :rules="rules" auto-label-width>
        <a-form-item field="name" label="角色名称">
          <a-input v-model="form.name" placeholder="如：运营人员" />
        </a-form-item>
        <a-form-item
          field="code"
          label="角色标识"
          extra="用于前端路由 meta.roles 权限匹配，创建后不可修改"
        >
          <a-input
            v-model="form.code"
            :disabled="!!editingId"
            placeholder="如：operator"
          />
        </a-form-item>
        <a-form-item field="description" label="描述">
          <a-textarea
            v-model="form.description"
            :max-length="255"
            show-word-limit
          />
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
  </div>
</template>

<script lang="ts" setup>
  import { reactive, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import type { FormInstance } from '@arco-design/web-vue/es/form';
  import dayjs from 'dayjs';
  import useLoading from '@/hooks/loading';
  import {
    queryRoles,
    createRole,
    updateRole,
    deleteRole,
    RoleRecord,
    RoleForm,
  } from '@/api/system';

  const { loading, setLoading } = useLoading(true);
  const list = ref<RoleRecord[]>([]);
  const keyword = ref('');
  const pagination = reactive({ current: 1, pageSize: 10, total: 0 });

  async function fetchData() {
    setLoading(true);
    try {
      const { data } = await queryRoles({
        current: pagination.current,
        pageSize: pagination.pageSize,
        name: keyword.value,
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
  function onPageChange(current: number) {
    pagination.current = current;
    fetchData();
  }

  const formRef = ref<FormInstance>();
  const formVisible = ref(false);
  const editingId = ref<number>();
  const emptyForm = (): RoleForm => ({
    name: '',
    code: '',
    description: '',
    status: 1,
  });
  const form = ref<RoleForm>(emptyForm());
  const rules = {
    name: [{ required: true, message: '请输入角色名称' }],
    code: [
      { required: true, message: '请输入角色标识' },
      { match: /^[a-zA-Z][\w-]*$/, message: '以字母开头，仅含字母数字 _ -' },
    ],
  };

  function openCreate() {
    editingId.value = undefined;
    form.value = emptyForm();
    formVisible.value = true;
  }
  function openEdit(r: RoleRecord) {
    editingId.value = r.id;
    form.value = {
      name: r.name,
      code: r.code,
      description: r.description,
      status: r.status,
    };
    formVisible.value = true;
  }
  async function onSubmit() {
    if (await formRef.value?.validate()) return false;
    try {
      if (editingId.value) await updateRole(editingId.value, form.value);
      else await createRole(form.value);
      Message.success('保存成功');
      fetchData();
      return true;
    } catch {
      return false;
    }
  }
  async function onDelete(r: RoleRecord) {
    await deleteRole(r.id);
    Message.success('已删除');
    fetchData();
  }

  fetchData();
</script>

<script lang="ts">
  export default {
    name: 'SystemRole',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
