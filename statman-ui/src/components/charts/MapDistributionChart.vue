<script>
import { Pie } from 'vue-chartjs'
export default {
  extends: Pie,
  props: {
    matchData: {
      type: Array
    }
  },
  data () {
    return {
      chartOptions: {
        legend: {
          labels: {
            fontColor: '#ebdbb2'
          }
        }
      }
    }
  },
  mounted () {
    var maps = []
    this.matchData.forEach(match => {
      maps.push(match.matchData.map)
    })
    const uniques = Array.from(new Set(maps))
    var values = []

    uniques.forEach(unique => {
      const count = maps.filter(x => x == unique).length
      values.push(count)
    })

    // Overwriting base render method with actual data.
    this.renderChart({
      labels: uniques,
      datasets: [{
          label: 'Times Played',
          backgroundColor: ['#d79921', '#98971a', '#cc241d', '#458588', '#689d6a'],
          data: values
        }],
    },
    this.chartOptions
    )
  }
}
</script>