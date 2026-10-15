<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.log']" />
    <a-card class="general-card" :title="$t('menu.system.log')">
      <a-row :gutter="16" class="toolbar">
        <a-col :flex="1">
          <a-space wrap :size="16">
            <a-input
              v-model="query.keyword"
              :placeholder="$t('system.log.keyword.placeholder')"
              allow-clear
              style="width: 240px"
              @press-enter="search"
            >
              <template #prefix><icon-search /></template>
            </a-input>
            <a-select
              v-model="query.method"
              :placeholder="$t('system.log.method')"
              allow-clear
              style="width: 140px"
            >
              <a-option v-for="m in methods" :key="m" :value="m">
                {{ m }}
              </a-option>
            </a-select>
            <a-radio-group v-model="query.result" type="button">
              <a-radio value="">{{ $t('system.common.all') }}</a-radio>
              <a-radio value="success">
                {{ $t('system.log.result.success') }}
              </a-radio>
              <a-radio value="fail">{{ $t('system.log.result.fail') }}</a-radio>
            </a-radio-group>
          </a-space>
        </a-col>
        <a-col flex="none">
          <a-space>
            <a-button type="primary" @click="search">
              {{ $t('system.common.search') }}
            </a-button>
            <a-button @click="reset">{{ $t('system.common.reset') }}</a-button>
          </a-space>
        </a-col>
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
          <a-table-column :title="$t('system.log.time')" :width="180">
            <template #cell="{ record }">
              {{ dayjs(record.createdAt).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.log.user')" :width="140">
            <template #cell="{ record }">
              <span v-if="record.username">{{ record.username }}</span>
              <span v-else class="muted">
                {{ $t('system.log.anonymous') }}
              </span>
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.log.path')">
            <template #cell="{ record }">
              <a-tag :color="methodColor[record.method]" size="small">
                {{ record.method }}
              </a-tag>
              <span class="path">{{ record.path }}</span>
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.log.result')" :width="120">
            <template #cell="{ record }">
              <a-badge
                :status="record.code === 20000 ? 'success' : 'danger'"
                :text="
                  record.code === 20000
                    ? $t('system.log.result.success')
                    : `${$t('system.log.result.fail')} · ${record.code}`
                "
              />
            </template>
          </a-table-column>
          <a-table-column :title="$t('system.log.latency')" :width="100">
            <template #cell="{ record }">{{ record.latencyMs }} ms</template>
          </a-table-column>
          <a-table-column
            :title="$t('system.log.ip')"
            data-index="ip"
            :width="140"
          />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import dayjs from 'dayjs';
  import useLoading from '@/hooks/loading';
  import { listLogs, OperationLog } from '@/api/system';

  const { loading, setLoading } = useLoading();
  const methods = ['POST', 'PUT', 'DELETE', 'PATCH'];
  const methodColor: Record<string, string> = {
    POST: 'green',
    PUT: 'arcoblue',
    DELETE: 'red',
    PATCH: 'orange',
  };
  const rows = ref<OperationLog[]>([]);
  const query = reactive({
    keyword: '',
    method: undefined as string | undefined,
    result: '',
  });
  const pagination = reactive({
    current: 1,
    pageSize: 20,
    total: 0,
    showTotal: true,
    showPageSize: true,
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const { data } = await listLogs({
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
  const search = () => {
    pagination.current = 1;
    fetchData();
  };
  const reset = () => {
    Object.assign(query, { keyword: '', method: undefined, result: '' });
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

  onMounted(fetchData);
</script>

<script lang="ts">
  export default {
    name: 'SystemLog',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px;
  }

  .toolbar {
    margin-bottom: 16px;
  }

  .path {
    margin-left: 8px;
    color: var(--color-text-2);
    font-family: var(--font-family-mono, monospace);
    font-size: 13px;
  }

  .muted {
    color: var(--color-text-4);
  }
</style>
