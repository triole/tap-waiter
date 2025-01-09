import { on } from 'mokapi'
import { fake } from 'mokapi/faker'

// content definitions
var category_type = ['software_recipe', 'container_build', 'document']
var purpose_goal = ['analysis', 'combine']
var purpose_datatype = ['s3', 'csv', 'hdf5', 'xroot']
var purpose_method = ['MCMC', 'correlation', 'likelihood', 'crossmatch']
var purpose_communities = ['astrophysics']
var license_name = ['Apache 2.0', 'GPL', 'CC0']

function rand(max) {
  return Math.floor(Math.random() * max)
}

function rand_range(min, max) {
  return Math.floor(Math.random() * (max - min) + min)
}

function pickrand(arr) {
  return arr[rand(arr.length)]
}

function pickrand_mult(arr, no) {
  var r = []
  for (var i = 0; i < no; i++) {
    r.push(pickrand(arr))
  }
  return r
}

// generators
function make_orcid() {
  return rand_range(100000000, 999999999)
}

function make_creator() {
  var r = {}
  r.name = fake({ type: 'string', format: '{name}' })
  r.orcid = make_orcid()
  return r
}

function make_person() {
  var r = {}
  r.name = fake({ type: 'string', format: '{name}' })
  r.orcid = make_orcid()
  r.mail = fake({ type: 'string', format: '{email}' })
  return r
}

function make_license() {
  var r = {}
  r.name = pickrand(license_name)
  r.url = fake({ type: 'string', format: '{url}' })
  return r
}

function make_repo() {
  var r = {}
  r.id = rand_range(1, 9000)
  r.name = fake({ type: 'string', format: '{productname}' }).toLowerCase().replace(/ /g, '-')
  r.readme_url = 'http://localhost:4466/punch/public/README.md'
  return r
}

// make iterators
function make_persons() {
  var r = []
  for (var i = 0; i < rand_range(1, 5); i++) {
    r.push(make_person())
  }
  return r
}

export default function () {
  on('http', function (request, response) {
    console.log(request)
    console.log(request.key)
    if (request.key === '/repos') {
      var d = []
      for (var i = 0; i < 5; i++) {
        d.push(make_repo())
      }
    } else {
      var d = {}
      d.category = {}
      d.category.type = pickrand(category_type)

      d.purpose = {}
      d.purpose.goal = pickrand(purpose_goal)
      d.purpose.datatype = pickrand(purpose_datatype)
      d.purpose.method = pickrand(purpose_method)
      d.purpose.communities = pickrand(purpose_communities)

      d.creators = make_persons()
      d.contact = make_persons()
      d.license = make_license()
    }
    response.data = d
  })
}
