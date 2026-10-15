<template>
  <div class="container">
    <Breadcrumb :items="['menu.system', 'menu.system.overview']" />
    <a-spin :loading="loading" style="width: 100%">
      <div class="hero">
        <div>
          <div class="hero-title">你好，{{ userStore.name }}</div>
          <div class="hero-desc">
            欢迎回到 Go Arco Admin · 今天是 {{ today }}
          </div>
        </div>
        <a-space>
          <a-button
            type="primary"
            @click="$router.push({ name: 'SystemUser' })"
          >
            <template #icon><icon-user-add /></template>管理用户
          </a-button>
          <a-button @click="$router.push({ name: 'SystemLog' })">
            查看日志
          </a-button>
        </a-space>
      </div>

      <a-grid :cols="{ xs: 1, sm: 2, lg: 4 }" :col-gap="16" :row-gap="16">
        <a-grid-item v-for="s in stats" :key="s.label">
          <div class="stat" :style="{ '--accent': s.color }">
            <div class="stat-icon"><component :is="s.icon" /></div>
            <a-statistic
              :title="s.label"
              :value="s.value"
              :value-from="0"
              animation
              show-group-separator
            />
          </div>
        </a-grid-item>
      </a-grid>

      <a-grid
        :cols="{ xs: 1, lg: 3 }"
        :col-gap="16"
        :row-gap="16"
        style="margin-top: 16px"
      >
        <a-grid-item :span="{ xs: 1, lg: 2 }">
          <a-card class="general-card" title="近 7 日操作趋势">
            <Chart height="300px" :options="chartOption" />
          </a-card>
        </a-grid-item>
        <a-grid-item>
          <a-card class="general-card" title="最新动态">
            <template #extra>
              <a-link @click="$router.push({ name: 'SystemLog' })">全部</a-link>
            </template>
            <a-timeline v-if="data?.recentLogs.length">
              <a-timeline-item
                v-for="log in data.recentLogs"
                :key="log.id"
                :dot-color="log.success ? '#00B42A' : '#F53F3F'"
              >
                <div class="tl-title">
                  <b>{{ log.username || '匿名' }}</b> {{ log.action }}
                </div>
                <div class="tl-time">{{ fromNow(log.createdAt) }}</div>
              </a-timeline-item>
            </a-timeline>
            <a-empty v-else />
          </a-card>
        </a-grid-item>
      </a-grid>
    </a-spin>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import dayjs from 'dayjs';
  import relativeTime from 'dayjs/plugin/relativeTime';
  import 'dayjs/locale/zh-cn';
  import { graphic } from 'echarts';
  import {
    IconUser,
    IconCheckCircle,
    IconSafe,
    IconHistory,
  } from '@arco-design/web-vue/es/icon';
  import useLoading from '@/hooks/loading';
  import useChartOption from '@/hooks/chart-option';
  import { useUserStore } from '@/store';
  import { getOverview, Overview } from '@/api/system';

  dayjs.extend(relativeTime);
  dayjs.locale('zh-cn');

  const userStore = useUserStore();
  const { loading, setLoading } = useLoading(true);
  const data = ref<Overview>();
  const today = dayjs().format('YYYY年M月D日 dddd');
  const fromNow = (v: string) => dayjs(v).fromNow();

  const stats = computed(() => [
    {
      label: '用户总数',
      value: data.value?.userCount ?? 0,
      icon: IconUser,
      color: '22, 93, 255',
    },
    {
      label: '启用账号',
      value: data.value?.activeCount ?? 0,
      icon: IconCheckCircle,
      color: '0, 180, 42',
    },
    {
      label: '角色数量',
      value: data.value?.roleCount ?? 0,
      icon: IconSafe,
      color: '114, 46, 209',
    },
    {
      label: '今日操作',
      value: data.value?.logsToday ?? 0,
      icon: IconHistory,
      color: '247, 114, 52',
    },
  ]);

  const { chartOption } = useChartOption((isDark) => ({
    grid: { left: 40, right: 20, top: 20, bottom: 30 },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: data.value?.trend.map((t) => t.day) ?? [],
      axisLine: { lineStyle: { color: isDark ? '#3f3f3f' : '#E5E8EF' } },
      axisLabel: { color: '#86909C' },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: {
        lineStyle: { type: 'dashed', color: isDark ? '#2E2E30' : '#E5E8EF' },
      },
      axisLabel: { color: '#86909C' },
    },
    series: [
      {
        name: '操作次数',
        type: 'line',
        smooth: true,
        symbolSize: 8,
        data: data.value?.trend.map((t) => t.count) ?? [],
        lineStyle: { width: 3, color: '#165DFF' },
        itemStyle: { color: '#165DFF', borderColor: '#fff', borderWidth: 2 },
        areaStyle: {
          color: new graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(22, 93, 255, 0.28)' },
            { offset: 1, color: 'rgba(22, 93, 255, 0)' },
          ]),
        },
      },
    ],
  }));

  getOverview()
    .then((res) => {
      data.value = res.data;
    })
    .finally(() => setLoading(false));
</script>

<script lang="ts">
  export default {
    name: 'SystemOverview',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
  .hero {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    padding: 28px 32px;
    color: #fff;
    background: linear-gradient(120deg, #165dff 0%, #4080ff 55%, #722ed1 100%);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(22, 93, 255, 0.18);
    &-title {
      font-weight: 600;
      font-size: 22px;
    }
    &-desc {
      margin-top: 6px;
      opacity: 0.85;
    }
    :deep(.arco-btn-secondary) {
      color: #fff;
      background: rgba(255, 255, 255, 0.16);
      &:hover {
        background: rgba(255, 255, 255, 0.26);
      }
    }
    :deep(.arco-btn-primary) {
      color: #165dff;
      background: #fff;
    }
  }
  .stat {
    display: flex;
    gap: 16px;
    align-items: center;
    padding: 20px;
    background: var(--color-bg-2);
    border-radius: 8px;
    &-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 48px;
      height: 48px;
      color: rgb(var(--accent));
      font-size: 24px;
      background: rgba(var(--accent), 0.1);
      border-radius: 12px;
    }
  }
  .tl-title {
    color: var(--color-text-1);
  }
  .tl-time {
    margin-top: 2px;
    color: var(--color-text-3);
    font-size: 12px;
  }
</style>
