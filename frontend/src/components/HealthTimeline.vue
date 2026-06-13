<script setup lang="ts">
import type { HealthRecord } from '../types'
import { healthApi } from '../api/health'

const props = defineProps<{
  records: HealthRecord[]
}>()

const emit = defineEmits<{
  refresh: []
}>()

const typeConfig: Record<string, { label: string; icon: string; color: string }> = {
  vaccine: { label: '疫苗', icon: '💉', color: '#3b82f6' },
  deworming: { label: '驱虫', icon: '🛡️', color: '#22c55e' },
  checkup: { label: '体检', icon: '🏥', color: '#f59e0b' },
}

function getConfig(type: string) {
  return typeConfig[type] || { label: type, icon: '📋', color: '#a8a29e' }
}

async function deleteRecord(record: HealthRecord) {
  if (!confirm(`确定删除 "${record.name}" 吗？`)) return
  try {
    await healthApi.delete(record.pet_id, record.id)
    emit('refresh')
  } catch {}
}
</script>

<template>
  <div class="health-timeline">
    <div v-if="records.length === 0" class="empty-state">
      <p>暂无健康记录</p>
    </div>
    <div v-for="record in records" :key="record.id" class="timeline-item">
      <div class="timeline-dot" :style="{ background: getConfig(record.type).color }"></div>
      <div class="timeline-content">
        <div class="timeline-header">
          <span class="timeline-icon">{{ getConfig(record.type).icon }}</span>
          <span class="timeline-type">{{ getConfig(record.type).label }}</span>
          <span class="timeline-name">{{ record.name }}</span>
          <button class="btn-icon" @click="deleteRecord(record)">×</button>
        </div>
        <div class="timeline-date">{{ record.date }}</div>
        <div v-if="record.next_date" class="timeline-next">
          下次：{{ record.next_date }}
        </div>
        <div v-if="record.report_url" class="timeline-report">
          <a :href="record.report_url" target="_blank" @click.stop>📄 查看体检报告</a>
        </div>
        <div v-if="record.note" class="timeline-note">{{ record.note }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.health-timeline {
  position: relative;
  padding-left: 1.5rem;
}

.health-timeline::before {
  content: '';
  position: absolute;
  left: 0.375rem;
  top: 0;
  bottom: 0;
  width: 1.5px;
  background: var(--bg-muted);
}

.empty-state {
  padding: 2rem;
  text-align: center;
  color: var(--fg-tertiary);
  font-size: var(--text-sm);
}

.timeline-item {
  position: relative;
  padding-bottom: 1.25rem;
}

.timeline-dot {
  position: absolute;
  left: -1.5rem;
  top: 0.375rem;
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  border: 2px solid var(--bg-elevated);
}

.timeline-content {
  background: var(--bg-subtle);
  border-radius: 0.625rem;
  padding: 0.75rem;
}

.timeline-header {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.timeline-icon { font-size: 0.875rem; }

.timeline-type {
  font-size: var(--text-xs);
  font-weight: 600;
  padding: 0.125rem 0.5rem;
  border-radius: 0.375rem;
  background: var(--bg-muted);
  color: var(--fg-secondary);
}

.timeline-name {
  flex: 1;
  font-weight: 600;
  font-size: var(--text-sm);
}

.timeline-date {
  font-size: var(--text-xs);
  color: var(--fg-tertiary);
  margin-top: 0.25rem;
}

.timeline-next {
  font-size: var(--text-xs);
  color: var(--warning);
  font-weight: 500;
  margin-top: 0.25rem;
}

.timeline-note {
  font-size: var(--text-xs);
  color: var(--fg-tertiary);
  margin-top: 0.25rem;
}

.timeline-report {
  margin-top: 0.25rem;
}

.timeline-report a {
  font-size: var(--text-xs);
  color: var(--info);
  text-decoration: none;
  font-weight: 500;
}

.timeline-report a:hover {
  text-decoration: underline;
}

.btn-icon {
  width: 1.5rem;
  height: 1.5rem;
  border: none;
  background: none;
  color: var(--fg-tertiary);
  cursor: pointer;
  border-radius: 0.375rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all var(--duration-fast);
}
.btn-icon:hover { background: var(--danger-soft); color: var(--danger); }
</style>
