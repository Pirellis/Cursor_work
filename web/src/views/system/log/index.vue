<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.log']" />
    <a-card class="general-card" :title="$t('menu.system.log')">
      <a-space style="margin-bottom: 16px">
        <a-input
          v-model="query.username"
          placeholder="操作人"
          allow-clear
          style="width: 200px"
          @press-enter="search"
        />
        <a-input
          v-model="query.action"
          placeholder="操作类型，如：登录"
          allow-clear
          style="width: 220px"
          @press-enter="search"
        />
        <a-button type="primary" @click="search">
          <template #icon><icon-search /></template>查询
        </a-button>
      </a-space>
      <a-table
        row-key="id"
        :loading="loading"
        :data="list"
        :pagination="pagination"
        :scroll="{ x: 1000 }"
        :bordered="false"
        @page-change="onPageChange"
      >
        <template #columns>
          <a-table-column title="时间" data-index="createdAt" :width="180">
            <template #cell="{ record }">
              {{ dayjs(record.createdAt).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
          </a-table-column>
          <a-table-column title="操作人" data-index="username" />
          <a-table-column title="操作" data-index="action" />
          <a-table-column title="请求" data-index="path">
            <template #cell="{ record }">
              <a-tag :color="methodColor[record.method] || 'gray'" size="small">
                {{ record.method }}
              </a-tag>
              <span class="path">{{ record.path }}</span>
            </template>
          </a-table-column>
          <a-table-column title="IP" data-index="ip" />
          <a-table-column title="耗时" data-index="latency" :width="90">
            <template #cell="{ record }">{{ record.latency }} ms</template>
          </a-table-column>
          <a-table-column title="结果" data-index="success" :width="90">
            <template #cell="{ record }">
              <a-badge
                :status="record.success ? 'success' : 'danger'"
                :text="record.success ? '成功' : '失败'"
              />
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { reactive, ref } from 'vue';
  import dayjs from 'dayjs';
  import useLoading from '@/hooks/loading';
  import { queryLogs, LogRecord } from '@/api/system';

  const { loading, setLoading } = useLoading(true);
  const list = ref<LogRecord[]>([]);
  const query = reactive({ username: '', action: '' });
  const pagination = reactive({ current: 1, pageSize: 15, total: 0 });
  const methodColor: Record<string, string> = {
    POST: 'arcoblue',
    PUT: 'orangered',
    DELETE: 'red',
  };

  async function fetchData() {
    setLoading(true);
    try {
      const { data } = await queryLogs({
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
  function onPageChange(current: number) {
    pagination.current = current;
    fetchData();
  }

  fetchData();
</script>

<script lang="ts">
  export default {
    name: 'SystemLog',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
  .path {
    margin-left: 8px;
    color: var(--color-text-2);
    font-family: Menlo, Consolas, monospace;
    font-size: 12px;
  }
</style>
