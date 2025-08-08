const express = require('express');
const swaggerJsdoc = require('swagger-jsdoc');
const swaggerUi = require('swagger-ui-express');

const app = express();
const PORT = 3000;

// Sample repository data
const repositories = [
  { id: 1, name: 'awesome-project', owner: 'user1', stars: 42 },
  { id: 2, name: 'cool-library', owner: 'user2', stars: 105 },
  { id: 3, name: 'utility-tool', owner: 'user3', stars: 27 }
];

// Swagger definition
const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Repository API',
      version: '1.0.0',
      description: 'API for managing code repositories',
    },
    servers: [
      {
        url: `http://localhost:${PORT}`,
      },
    ],
    components: {
      schemas: {
        Repository: {
          type: 'object',
          properties: {
            id: {
              type: 'integer',
              example: 1
            },
            name: {
              type: 'string',
              example: 'awesome-project'
            },
            owner: {
              type: 'string',
              example: 'user1'
            },
            stars: {
              type: 'integer',
              example: 42
            }
          }
        }
      }
    }
  },
  apis: [], // No external files for this example
};

const specs = swaggerJsdoc(options);

// API routes
/**
 * @swagger
 * /repositories:
 *   get:
 *     summary: Get all repositories
 *     responses:
 *       200:
 *         description: A list of repositories
 *         content:
 *           application/json:
 *             schema:
 *               type: array
 *               items:
 *                 $ref: '#/components/schemas/Repository'
 */
app.get('/repositories', (req, res) => {
  res.json(repositories);
});

/**
 * @swagger
 * /repositories/{id}:
 *   get:
 *     summary: Get a repository by ID
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     responses:
 *       200:
 *         description: A single repository
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Repository'
 *       404:
 *         description: Repository not found
 */
app.get('/repositories/:id', (req, res) => {
  const repo = repositories.find(r => r.id === parseInt(req.params.id));
  if (!repo) return res.status(404).send('Repository not found');
  res.json(repo);
});

// Serve Swagger UI
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(specs));

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
  console.log(`Swagger UI available at http://localhost:${PORT}/api-docs`);
});