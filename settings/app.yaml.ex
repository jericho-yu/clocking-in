# 月度打卡
checking-in:
  filename: ./考勤5月-智冶互联.xlsx
  start-row: 5
  # 打卡时间相关配置
  clock-in-time:
    sheet-name: 打卡时间
    overtimes: H,L,M,N.O.P,W,AC,AD,AJ
  # 月度汇总相关配置
  month:
    sheet-name: 月度汇总 (2)
    overtimes: H,L,M,N.O.P,W,AC,AD,AJ
  
# 月度汇总
collect: 
  filename: ./加班、调休统计表-2024-智冶互联.xlsx
  start-row: 5