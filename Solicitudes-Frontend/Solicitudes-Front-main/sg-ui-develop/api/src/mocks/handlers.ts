import { http, HttpResponse } from "msw";
import jwt from "jsonwebtoken";

import { configService } from "../configService";

const mockAddresses = [
  {
    address_street: "HEROES DE LAS MALVINAS",
    address_number: 0,
    abl_number: 113028,
    property_id: 1,
  },
  {
    address_street: "HEROES DE LAS MALVINAS",
    address_number: 357,
    abl_number: 128663,
    property_id: 2,
  },
  {
    address_street: "SAN MARTIN",
    address_number: 123,
    abl_number: 456789,
    property_id: 3,
  },
  {
    address_street: "AV. BELGRANO",
    address_number: 50,
    abl_number: 987654,
    property_id: 4,
  },
];

const mockActivities = [
  {
    id: 1,
    description: "Ejecutar solados (piso)",
  },
  {
    id: 2,
    description: "Cambiar revestimientos",
  },
  {
    id: 3,
    description: "Terraplenar y rellenar terrenos",
  },
  {
    id: 4,
    description: "Cambiar el material de cubierta de techos",
  },
  {
    id: 5,
    description: "Ejecutar cielorrasos",
  },
  {
    id: 6,
    description: "Revocar cercas al frente",
  },
  {
    id: 7,
    description: "Ejecutar revoques exteriores o trabajos similares",
  },
  {
    id: 8,
    description: "Limpiar o pintar las fachadas principales",
  },
];

const inMemoryDb: any[] = [];
const usersMemoryDb: any[] = [
  {
    cuil: "20348134664",
    dni: "34813466",
    first_name: "LESLIE",
    last_name: "ANN CHRISTINE",
    email: "sasad@gmail.com",
    phone: "1212121212",
    email_validated: true,
  },
];

export const handlers = [
  http.get(
    configService.baseRequestApi + "/requests/address/autocomplete",
    ({ request }) => {
      const url = new URL(request.url);
      const query = url.searchParams.get("q")?.toLowerCase();
      if (!query) {
        return new HttpResponse(
          JSON.stringify({ message: "Query no proporcionado" }),
          {
            status: 400,
            headers: {
              "Content-Type": "application/json",
            },
          }
        );
      }

      const filteredAddresses = mockAddresses.filter((address) =>
        address.address_street.toLowerCase().includes(query)
      );

      return new Response(JSON.stringify({ suggestions: filteredAddresses }), {
        headers: {
          "Content-Type": "application/json",
        },
      });
    }
  ),
  http.get(configService.baseRequestApi + "/site/activities", () => {
    return new Response(JSON.stringify(mockActivities), {
      headers: {
        "Content-Type": "application/json",
      },
    });
  }),
  http.get(
    configService.baseRequestApi + "/requests/protected/get/all/cuil",
    () => {
      return new Response(JSON.stringify({ requests: inMemoryDb }), {
        headers: {
          "Content-Type": "application/json",
        },
      });
    }
  ),
  http.post(
    configService.baseRequestApi + "/requests/protected/create",
    async ({ request }) => {
      const body = await request.json();

      const newRequest = {
        body,
        id: `ABCDE-${Math.random().toString(36).substr(2, 9)}`,
        type: "Aviso de Obra",
        created_at: new Date().toISOString(),
        status: "inProgress",
      };

      inMemoryDb.push(newRequest);

      return new Response(
        JSON.stringify({ success: true, code: "ABCDE-12345-FGHI-6789" }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }
  ),
  http.get(
    configService.baseLoginApi + "/auth/protected/login",
    ({ request }) => {
      const url = new URL(request.url);
      const cuil = url.searchParams.get("cuil")?.toLowerCase();
      if (!cuil) {
        return new HttpResponse(
          JSON.stringify({ message: "Cuil no proporcionado" }),
          {
            status: 400,
            headers: {
              "Content-Type": "application/json",
            },
          }
        );
      }

      const token = generateToken(cuil);
      if (token === "") {
        return new HttpResponse(JSON.stringify({ message: "Invalid token" }), {
          status: 400,
          headers: {
            "Content-Type": "application/json",
          },
        });
      }

      return new Response(
        JSON.stringify({
          message: "Login successful",
          cuil,
          access_token: token,
        }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }
  ),
  http.get(configService.baseLoginApi + "/auth/protected/hi", ({ request }) => {
    const authHeader = request.headers.get("authorization");
    if (!authHeader) {
      return new HttpResponse(JSON.stringify({ message: "Invalid token" }), {
        status: 401,
        headers: {
          "Content-Type": "application/json",
        },
      });
    }

    const token = authHeader.split(" ")[1];
    if (!token) {
      return new HttpResponse(JSON.stringify({ message: "Invalid token" }), {
        status: 401,
        headers: {
          "Content-Type": "application/json",
        },
      });
    }

    const secretKey = process.env.JWT_SECRET;
    if (!secretKey) {
      return new HttpResponse(
        JSON.stringify({ message: "Internal Server Error" }),
        {
          status: 500,
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    }

    try {
      jwt.verify(token, secretKey);
      return new Response(
        JSON.stringify({
          message: "Login successful",
          access_token: token,
        }),
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );
    } catch (err) {
      return new HttpResponse(JSON.stringify({ message: "Invalid token" }), {
        status: 401,
        headers: {
          "Content-Type": "application/json",
        },
      });
    }
  }),
  http.post(configService.baseUsersApi + "/users", async ({ request }) => {
    const body = await request.json();
    // cuil: "20348134664",
    // dni: "34813466",
    // first_name: "LESLIE",
    // last_name: "ANN CHRISTINE",
    // email: "sasad@gmail.com",
    // phone: "1212121212",
    // email_validated: true,

    // const { cuil, dni, first_name, last_name, email, phone } = body;

    // const newUser = {
    //   cuil,
    //   type: "Aviso de Obra",
    //   created_at: new Date().toISOString(),
    //   status: "pending",
    // };

    usersMemoryDb.push(body);

    return new Response(JSON.stringify(body), {
      headers: {
        "Content-Type": "application/json",
      },
    });
  }),
  http.get(configService.baseUsersApi + "/users", ({ request }) => {
    const url = new URL(request.url);
    const cuil = url.searchParams.get("cuil")?.toLowerCase();

    // const user = usersMemoryDb.find((u) => u.cuil.toLowerCase() === cuil);

    // if (!user) {
    //   return new Response(JSON.stringify({ error: "Usuario no encontrado" }), {
    //     status: 404,
    //     headers: {
    //       "Content-Type": "application/json",
    //     },
    //   });
    // }

    return new Response(
      JSON.stringify({
        cuil: cuil,
        dni: "34813466",
        first_name: "LESLIE",
        last_name: "ANN CHRISTINE",
        email: "sasad@gmail.com",
        phone: "1212121212",
        email_validated: true,
      }),
      {
        headers: {
          "Content-Type": "application/json",
        },
      }
    );
  }),
];

const generateToken = (cuil: string): string => {
  const claims = {
    sub: cuil,
    exp: Math.floor(Date.now() / 1000) + 24 * 60 * 60, // 24 horas
    iat: Math.floor(Date.now() / 1000), // Tiempo actual
  };

  const secretKey = process.env.JWT_SECRET;
  if (!secretKey) {
    return "";
  }

  return jwt.sign(claims, secretKey, { algorithm: "HS256" });
};
