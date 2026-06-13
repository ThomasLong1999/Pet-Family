<script setup lang="ts">
import { computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import type { WeightRecord } from '../types'
import { t } from '../composables/useI18n'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const props = defineProps<{
  weights: WeightRecord[]
}>()

// Design tokens as hex literals (ECharts can't read CSS variables directly)
const FG_PRIMARY = '#1c1917'
const FG_TERTIARY = '#78716c' // matches --fg-tertiary (WCAG AA)
const BG_SUBTLE = '#f5f5f4'
const BG_MUTED = '#e7e5e4'

const chartOption = computed(() => {
  const sorted = [...props.weights].sort((a, b) =>
    a.recorded_at.localeCompare(b.recorded_at)
  )

  return {
    grid: { left: 45, right: 16, top: 16, bottom: 30 },
    tooltip: {
      trigger: 'axis' as const,
      formatter: (params: { axisValue: string; value: number }[]) => {
        const p = params[0]
        return `${p.axisValue}<br/><strong>${p.value} kg</strong>`
      },
      backgroundColor: 'rgba(255,255,255,0.95)',
      borderColor: BG_MUTED,
      textStyle: { color: FG_PRIMARY, fontSize: 13 },
    },
    xAxis: {
      type: 'category' as const,
      data: sorted.map(w => w.recorded_at.slice(0, 10)),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: FG_TERTIARY, fontSize: 11 },
    },
    yAxis: {
      type: 'value' as const,
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: BG_SUBTLE } },
      axisLabel: { color: FG_TERTIARY, fontSize: 11, formatter: '{value} kg' },
    },
    series: [{
      type: 'line',
      data: sorted.map(w => w.weight),
      smooth: true,
      symbol: 'circle',
      symbolSize: 6,
      lineStyle: { width: 2.5, color: FG_PRIMARY },
      itemStyle: { color: FG_PRIMARY, borderColor: '#fff', borderWidth: 2 },
      areaStyle: {
        color: {
          type: 'linear' as const,
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(28, 25, 23, 0.08)' },
            { offset: 1, color: 'rgba(28, 25, 23, 0)' },
          ],
        },
      },
    }],
  }
})
</script>

<template>
  <div class="weight-chart">
    <div v-if="weights.length === 0" class="empty-chart">
      <p>{{ t('weight.empty') }}</p>
    </div>
    <VChart v-else :option="chartOption" autoresize style="height: 200px; width: 100%;" />
  </div>
</template>

<style scoped>
.weight-chart {
  background: var(--bg-subtle);
  border-radius: 0.75rem;
  padding: 0.5rem;
}

.empty-chart {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  color: var(--fg-tertiary);
  font-size: var(--text-sm);
}
</style>
