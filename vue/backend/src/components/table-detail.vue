<template>
  <el-descriptions :title="title" :column="column" border>
    <el-descriptions-item v-for="item in list" :key="item.prop" :span="item.span">
      <template #label> {{ item.label }} </template>
      <slot :name="item.prop" :row="row">
        {{ item.formatter? item.formatter(row) : item.value || row[item.prop] }}
      </slot>
    </el-descriptions-item>
  </el-descriptions>
</template>

<script lang="ts" setup>
const props = defineProps({
  data: {
    type: Object,
    required: true,
  }
});
const { row, title, column = 2, list } = props.data;
</script>

<style scoped>
/* 整体描述列表容器样式 */
.el-descriptions {
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

/* 描述项样式 */
.el-descriptions__item {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  border-bottom: 1px solid #e4e7ed;
  padding: 12px 0;
}

/* 描述项标签样式 */
.el-descriptions__item__label {
  flex: 0 0 auto;
  width: 150px;
  margin-right: 15px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}

/* 描述项内容样式 */
.el-descriptions__item__content {
  flex: 1 1 auto;
  max-height: 100px;
  overflow-y: auto;
  word-break: break-all;
}

/* 当描述项数量较多时，调整布局 */
@media (max - width: 600px) {
  .el-descriptions__item__label {
    width: 100%;
    margin-right: 0;
    margin-bottom: 5px;
  }
  .el-descriptions__item__content {
    width: 100%;
  }
}
</style>