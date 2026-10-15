<template>
  <a-grid :cols="{ xs: 1, sm: 2, lg: 4 }" :col-gap="16" :row-gap="16">
    <a-grid-item v-for="item in items" :key="item.key">
      <div class="overview-card" :style="{ '--accent': item.color }">
        <div class="overview-card-icon">
          <component :is="item.icon" />
        </div>
        <a-statistic
          :title="$t(`system.overview.${item.key}`)"
          :value="overview[item.key]"
          :value-from="0"
          animation
          show-group-separator
        />
      </div>
    </a-grid-item>
  </a-grid>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { getOverview, Overview } from '@/api/system';

  const overview = ref<Overview>({
    users: 0,
    activeUsers: 0,
    roles: 0,
    todayOperations: 0,
  });

  const items: { key: keyof Overview; icon: string; color: string }[] = [
    { key: 'users', icon: 'icon-user-group', color: '22, 93, 255' },
    { key: 'activeUsers', icon: 'icon-check-circle', color: '0, 180, 42' },
    { key: 'roles', icon: 'icon-safe', color: '114, 46, 209' },
    { key: 'todayOperations', icon: 'icon-thunderbolt', color: '255, 125, 0' },
  ];

  const load = async () => {
    const { data } = await getOverview();
    overview.value = data;
  };

  onMounted(load);
  defineExpose({ reload: load });
</script>

<style scoped lang="less">
  .overview-card {
    position: relative;
    display: flex;
    gap: 16px;
    align-items: center;
    padding: 20px;
    overflow: hidden;
    background: linear-gradient(
      135deg,
      rgba(var(--accent), 0.08) 0%,
      var(--color-bg-2) 70%
    );
    border: 1px solid rgba(var(--accent), 0.12);
    border-radius: 8px;
    transition: transform 0.2s ease, box-shadow 0.2s ease;

    &:hover {
      box-shadow: 0 8px 24px rgba(var(--accent), 0.12);
      transform: translateY(-2px);
    }

    &-icon {
      display: flex;
      flex-shrink: 0;
      align-items: center;
      justify-content: center;
      width: 48px;
      height: 48px;
      color: rgb(var(--accent));
      font-size: 24px;
      background: rgba(var(--accent), 0.12);
      border-radius: 12px;
    }

    :deep(.arco-statistic-title) {
      margin-bottom: 4px;
      color: var(--color-text-3);
      font-size: 13px;
    }

    :deep(.arco-statistic-value) {
      font-weight: 600;
      font-size: 26px;
    }
  }
</style>
