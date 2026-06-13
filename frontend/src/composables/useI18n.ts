import { ref, computed } from 'vue'

export type Locale = 'zh' | 'en'

const locale = ref<Locale>('zh')

export function setLocale(l: Locale) { locale.value = l }
export function getLocale() { locale.value
  return locale.value
}
export const isZh = computed(() => locale.value === 'zh')

const messages: Record<Locale, Record<string, string>> = {
  zh: {
    // Header
    'app.title': '自家毛孩',
    'btn.addPet': '添加宠物',
    'lang.toggle': 'Eng',

    // PetCard
    'pet.age.months': '{n}个月',
    'pet.age.years': '{y}岁{m}个月',
    'pet.age.yearsRound': '{y}岁',
    'pet.age.lessThan1': '< 1个月',
    'pet.weight': 'kg',
    'reminder.daysLeft': '{n}天后到期',

    // PetForm
    'form.addPet': '添加宠物',
    'form.editPet': '编辑宠物',
    'form.avatar': '头像',
    'form.addPhoto': '添加照片',
    'form.type': '类型',
    'form.name': '名字',
    'form.name.ph': '宠物名字',
    'form.breed': '品种',
    'form.breed.ph': '如：英短蓝猫',
    'form.gender': '性别',
    'form.gender.male': '公 ♂',
    'form.gender.female': '母 ♀',
    'form.birthday': '生日',
    'form.birthday.ph': '如 20240301',
    'form.adoptedAt': '领养日期',
    'form.adoptedAt.ph': '如 20240415',
    'form.passedAt': '去喵星日期',
    'form.passedAt.other': '离开日期',
    'form.passedAt.ph': '如 20260614',
    'form.color': '毛色',
    'form.color.ph': '如：蓝灰色',
    'form.note': '备注',
    'form.note.ph': '备注信息...',
    'form.save': '添加宠物',
    'form.saving': '保存中...',
    'form.cancel': '取消',
    'form.saveChanges': '保存',
    'form.delete': '删除宠物',

    // Species
    'species.cat': '🐱 猫',
    'species.dog': '🐶 狗',
    'species.hamster': '🐹 仓鼠',
    'species.rabbit': '🐰 兔子',

    // PetDetail tabs
    'tab.info': '档案',
    'tab.weight': '体重',
    'tab.health': '健康',
    'tab.photos': '相册',
    'detail.editInfo': '编辑信息',
    'detail.markPassed': '🌈 去喵星了',
    'detail.markPassed.other': '🌈 回快乐老家了',
    'detail.clearPassed': '取消标记',
    'detail.passedDate': '去喵星日期',
    'detail.passedDate.other': '离开日期',

    // Info labels
    'info.breed': '品种',
    'info.gender': '性别',
    'info.birthday': '生日',
    'info.color': '毛色',
    'info.adoptedAt': '领养日期',
    'info.note': '备注',

    // Weight
    'weight.title': '体重趋势',
    'weight.record': '记录体重',
    'weight.date': '日期',
    'weight.date.ph': '如 20260614',
    'weight.note': '备注',
    'weight.note.ph': '可选备注',
    'weight.empty': '暂无体重记录',
    'weight.kg': 'kg',

    // Health
    'health.title': '健康记录',
    'health.add': '添加记录',
    'health.type': '类型',
    'health.type.vaccine': '💉 疫苗',
    'health.type.deworming': '🛡️ 驱虫',
    'health.type.checkup': '🏥 体检',
    'health.quickSelect': '快捷选择',
    'health.name': '项目名称',
    'health.name.ph': '选择快捷项或手动输入',
    'health.date': '实施日期',
    'health.nextDate': '下次到期',
    'health.report': '体检报告',
    'health.report.ph': '点击上传报告文件',
    'health.note': '备注',
    'health.note.ph': '品牌、医院等',
    'health.empty': '暂无健康记录',
    'health.next': '下次：',
    'health.viewReport': '📄 查看体检报告',

    // Photos
    'photos.title': '成长相册',
    'photos.upload': '上传照片',
    'photos.empty': '暂无照片',
    'photos.dialog.title': '上传成长照片',
    'photos.dialog.select': '选择照片',
    'photos.ageGroup': '年龄分组',
    'photos.caption': '描述（可选）',
    'photos.caption.ph': '照片描述',
    'photos.confirm': '确认上传',
    'photos.uploading': '上传中...',
    'photos.ageGroup.month': '{n}个月',
    'photos.ageGroup.year': '{n}岁',

    // Avatar cropper
    'crop.title': '调整头像',
    'crop.confirm': '确认',
    'crop.hint': '拖动图片调整位置',

    // Dashboard
    'empty.title': '还没有宠物',
    'empty.desc': '点击上方按钮添加你的第一只宠物吧',

    // Passed prompt
    'passed.prompt': '请输入离开的日期（支持 20260614 或 2026-06-14）：',
    'passed.clear.confirm': '确定要取消{name}的标记吗？',
    'passed.delete.confirm': '确定要删除{name}吗？此操作不可恢复。',
  },
  en: {
    'app.title': 'Pet Family',
    'btn.addPet': 'Add Pet',
    'lang.toggle': '中',

    'pet.age.months': '{n}mo',
    'pet.age.years': '{y}y{m}mo',
    'pet.age.yearsRound': '{y}y',
    'pet.age.lessThan1': '< 1mo',
    'pet.weight': 'kg',
    'reminder.daysLeft': '{n}d left',

    'form.addPet': 'Add Pet',
    'form.editPet': 'Edit Pet',
    'form.avatar': 'Avatar',
    'form.addPhoto': 'Add Photo',
    'form.type': 'Type',
    'form.name': 'Name',
    'form.name.ph': 'Pet name',
    'form.breed': 'Breed',
    'form.breed.ph': 'e.g. British Shorthair',
    'form.gender': 'Gender',
    'form.gender.male': 'Male ♂',
    'form.gender.female': 'Female ♀',
    'form.birthday': 'Birthday',
    'form.birthday.ph': 'e.g. 20240301',
    'form.adoptedAt': 'Adopted Date',
    'form.adoptedAt.ph': 'e.g. 20240415',
    'form.passedAt': 'Passed Date',
    'form.passedAt.other': 'Passed Date',
    'form.passedAt.ph': 'e.g. 20260614',
    'form.color': 'Color',
    'form.color.ph': 'e.g. Blue-gray',
    'form.note': 'Note',
    'form.note.ph': 'Notes...',
    'form.save': 'Add Pet',
    'form.saving': 'Saving...',
    'form.cancel': 'Cancel',
    'form.saveChanges': 'Save',
    'form.delete': 'Delete Pet',

    'species.cat': '🐱 Cat',
    'species.dog': '🐶 Dog',
    'species.hamster': '🐹 Hamster',
    'species.rabbit': '🐰 Rabbit',

    'tab.info': 'Info',
    'tab.weight': 'Weight',
    'tab.health': 'Health',
    'tab.photos': 'Album',
    'detail.editInfo': 'Edit',
    'detail.markPassed': '🌈 Rainbow Bridge',
    'detail.markPassed.other': '🌈 Rainbow Bridge',
    'detail.clearPassed': 'Clear',
    'detail.passedDate': 'Passed Date',
    'detail.passedDate.other': 'Passed Date',

    'info.breed': 'Breed',
    'info.gender': 'Gender',
    'info.birthday': 'Birthday',
    'info.color': 'Color',
    'info.adoptedAt': 'Adopted',
    'info.note': 'Note',

    'weight.title': 'Weight Trend',
    'weight.record': 'Record Weight',
    'weight.date': 'Date',
    'weight.date.ph': 'e.g. 20260614',
    'weight.note': 'Note',
    'weight.note.ph': 'Optional note',
    'weight.empty': 'No weight records',
    'weight.kg': 'kg',

    'health.title': 'Health Records',
    'health.add': 'Add Record',
    'health.type': 'Type',
    'health.type.vaccine': '💉 Vaccine',
    'health.type.deworming': '🛡️ Deworming',
    'health.type.checkup': '🏥 Checkup',
    'health.quickSelect': 'Quick Select',
    'health.name': 'Name',
    'health.name.ph': 'Select or type',
    'health.date': 'Date',
    'health.nextDate': 'Next Due',
    'health.report': 'Report',
    'health.report.ph': 'Upload report',
    'health.note': 'Note',
    'health.note.ph': 'Brand, hospital...',
    'health.empty': 'No health records',
    'health.next': 'Next: ',
    'health.viewReport': '📄 View Report',

    'photos.title': 'Growth Album',
    'photos.upload': 'Upload',
    'photos.empty': 'No photos yet',
    'photos.dialog.title': 'Upload Photo',
    'photos.dialog.select': 'Select Photo',
    'photos.ageGroup': 'Age Group',
    'photos.caption': 'Caption (optional)',
    'photos.caption.ph': 'Photo caption',
    'photos.confirm': 'Upload',
    'photos.uploading': 'Uploading...',
    'photos.ageGroup.month': '{n}mo',
    'photos.ageGroup.year': '{n}y',

    'crop.title': 'Crop Avatar',
    'crop.confirm': 'Confirm',
    'crop.hint': 'Drag to reposition',

    'empty.title': 'No pets yet',
    'empty.desc': 'Click the button above to add your first pet',

    'passed.prompt': 'Enter the date (e.g. 20260614 or 2026-06-14):',
    'passed.clear.confirm': 'Remove the mark from {name}?',
    'passed.delete.confirm': 'Delete {name}? This cannot be undone.',
  },
}

export function t(key: string, params?: Record<string, string | number>): string {
  let text = messages[locale.value]?.[key] ?? messages.zh[key] ?? key
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      text = text.replace(`{${k}}`, String(v))
    }
  }
  return text
}

/** Species-aware key resolver */
export function speciesKey(species: string, key: string): string {
  if (species === 'cat') return t(key)
  return t(key + '.other', undefined as any) !== key + '.other' ? t(key + '.other', undefined as any) : t(key)
}
